package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"vincentinttsh/openclass-system/config"
	"vincentinttsh/openclass-system/model"
	"vincentinttsh/openclass-system/pkg/mode"
	"vincentinttsh/openclass-system/router"
	"vincentinttsh/openclass-system/view"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/django"
	"github.com/joho/godotenv"
)

var app *fiber.App
var serverConfig config.Config
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

		serverConfig.Production()
	} else {
		err := godotenv.Load()

		if err != nil {
			log.Panic("Error loading .env file")
		}

		engine.Debug(true)
		engine.Reload(true)

		serverConfig.Dev()
	}

	app = fiber.New(fiberConfig)

	model.InitFunc(&serverConfig)
	view.InitFunc(&serverConfig)
	router.SetupRouter(app, &serverConfig)
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

		if err := app.Listen(serverConfig.Address); err != nil {
			log.Panic(err)
		}

		cleanup()
	}
}
