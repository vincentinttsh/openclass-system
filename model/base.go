package model

import (
	"log"
	"os"
	"time"
	"vincentinttsh/openclass-system/pkg/mode"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func init() {
	var err error
	if mode.Mode() == mode.DebugMode {
		config := &gorm.Config{
			PrepareStmt: true,
		}
		if os.Getenv("DB_LOG") == "True" {
			config.Logger = logger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
				logger.Config{
					SlowThreshold:             time.Second, // Slow SQL threshold
					LogLevel:                  logger.Info, // Log level
					IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
					Colorful:                  true,        // Disable color
				},
			)
		}
		db, err = gorm.Open(sqlite.Open("dev.sqlite"), config)
	} else if mode.Mode() == mode.ReleaseMode {
		db, err = gorm.Open(postgres.Open(os.Getenv("DB_URL")), &gorm.Config{
			PrepareStmt: true,
		})
	} else {
		db, err = gorm.Open(sqlite.Open("test.sqlite"), &gorm.Config{})
	}

	if err != nil {
		panic("failed to connect database")
	}

	panicAtError(db.AutoMigrate(&Organization{}))
	panicAtError(db.AutoMigrate(&User{}))
	panicAtError(db.AutoMigrate(&GoogleOauth{}))
	panicAtError(db.AutoMigrate(&BaseClass{}))
}

func panicAtError(err error) {
	if err != nil {
		panic(err)
	}
}

// BaseModel base model
type BaseModel struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	ID        uint           `gorm:"not null;autoIncrement;primaryKey"`
}
