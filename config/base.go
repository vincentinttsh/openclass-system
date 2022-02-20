package config

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Config web server config
type Config struct {
	Address        string
	SessionStorage fiber.Storage
	DB             *gorm.DB
	Domain         string
}
