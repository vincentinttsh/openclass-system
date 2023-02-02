package model

import (
	"log"
	"os"
	"time"
	"vincentinttsh/openclass-system/internal/mode"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	dateFormat     string = "02/01/2006"
	pureTimeFormat string = "03:04 PM"
	timeFormat     string = "02/01/200603:04 PM"
)

var db *gorm.DB

// SQLBasePK => uint64
type SQLBasePK uint64

func init() {
	var err error
	config := &gorm.Config{
		PrepareStmt: true,
	}
	config.Logger = logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,          // Disable color
		},
	)
	switch mode.Mode() {
	case mode.DebugMode:
		if os.Getenv("DB_LOG") == "true" {
			config.Logger = config.Logger.LogMode(logger.Info)
		} else {
			config.Logger = config.Logger.LogMode(logger.Warn)
		}
		db, err = gorm.Open(sqlite.Open("dev.sqlite"), config)
	case mode.ReleaseMode:
		if os.Getenv("DB_LOG") == "true" {
			config.Logger = config.Logger.LogMode(logger.Info)
		}
		db, err = gorm.Open(postgres.Open(os.Getenv("DB_URL")), config)
	default:
		db, err = gorm.Open(sqlite.Open("test.sqlite"), config)
	}

	if err != nil {
		panic("failed to connect database")
	}

	panicAtError(db.AutoMigrate(&Organization{}))
	panicAtError(db.AutoMigrate(&User{}))
	panicAtError(db.AutoMigrate(&GoogleOauth{}))
	panicAtError(db.AutoMigrate(&Course{}))
	panicAtError(db.AutoMigrate(&CourseReservation{}))
	panicAtError(db.AutoMigrate(&SHDesignDetail{}))
	panicAtError(db.AutoMigrate(&SHDesign{}))
	panicAtError(db.AutoMigrate(&SHPreparation{}))
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
	ID        SQLBasePK      `gorm:"not null;autoIncrement;primaryKey"`
}

// DurationBaseModel base model for duration
type DurationBaseModel struct {
	Start     time.Time `gorm:"type:date"`
	End       time.Time `gorm:"type:date check (\"end\" > \"start\")"`
	Date      string    `gorm:"-:all" form:"date" i18n:"教學日期"`
	StartTime string    `gorm:"-:all" form:"startTime" i18n:"教學開始時間"`
	EndTime   string    `gorm:"-:all" form:"endTime" i18n:"教學結束時間"`
}

// DurationBaseInterface duration base interface
type DurationBaseInterface interface {
	GetTimeString() (string, string, string)
	SetStartTime(time.Time)
	SetEndTime(time.Time)
}
