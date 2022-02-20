package config

import (
	"github.com/gofiber/storage/memory"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Dev  web server development config
func (c *Config) Dev() {
	var err error

	c.Address = "localhost:8000"
	c.SessionStorage = memory.New()
	c.Domain = "localhost"
	c.DB, err = gorm.Open(sqlite.Open("dev.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
}
