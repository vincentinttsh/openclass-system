package config

import (
	"github.com/gofiber/storage/mongodb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Production  web server production config
func (c *Config) Production() {
	var err error

	c.Address = ":8000"
	c.SessionStorage = mongodb.New(mongodb.Config{
		Database: "fiber",
	})
	c.Domain = "localhost"
	c.DB, err = gorm.Open(postgres.Open(""), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
}
