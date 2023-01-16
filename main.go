package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"vincentinttsh/openclass-system/pkg/mode"
	"vincentinttsh/openclass-system/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/django"
)

var app *fiber.App
var engine *django.Engine

func start() {
	var fiberConfig fiber.Config

	fiberConfig.AppName = "OpenClass System"
	engine = django.New("./web/template", ".html")
	fiberConfig.Views = engine

	if mode.Mode() == mode.ReleaseMode {
		fiberConfig.Prefork = true
		fiberConfig.DisableStartupMessage = true
		fiberConfig.ReadTimeout = 10 * time.Second
	} else {
		engine.Reload(true)
	}

	app = fiber.New(fiberConfig)

	app.Use(favicon.New(favicon.Config{
		File: "./web/static/favicon.ico",
	}))

	if mode.Mode() == mode.ReleaseMode {
		app.Static("/static", "./web/static", fiber.Static{
			Compress: true,
			MaxAge:   31536000,
		})
	} else {
		app.Static("/static", "./web/static", fiber.Static{
			Browse:        true,
			CacheDuration: 0 * time.Second,
		})
	}

	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${latency} ${method} ${path}\n",
		TimeFormat: "2006-Jan-02 15:04:05",
		TimeZone:   "Asia/Taipei",
	}))

	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
	router.SetupRouter(app)
}

func cleanup() {
	fmt.Println("Running cleanup tasks...")
}

func main() {
	arg := ""
	if len(os.Args) > 1 {
		arg = os.Args[1]
	}
	if arg == "ping" {
		resp, err := http.Get("http://localhost:8080/ping")
		if err != nil {
			os.Exit(1)
		}
		if resp.StatusCode != http.StatusOK {
			os.Exit(1)
		}
		os.Exit(0)
	} else {
		start()

		// Wait for interrupt signal to gracefully shutdown the server with
		// a timeout of 5 seconds.
		quit := make(chan os.Signal, 1)
		// kill (no param) default send syscall.SIGTERM
		// kill -2 is syscall.SIGINT
		// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			<-quit
			fmt.Println("Gracefully shutting down...")
			if err := app.Shutdown(); err != nil {
				fmt.Println("Error in app.Shutdown():", err)
			}
		}()

		if err := app.Listen(os.Getenv("PORT")); err != nil {
			log.Panic(err)
		}

		cleanup()
	}
}
