package database

import (
	"sync"

	"github.com/pnaskardev/URL-Shortner-V1/shortener-service/infrastructure/database/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var postgresClient *gorm.DB = nil
var PostgreSyncOnce sync.Once

func ConnectToPostgres() *gorm.DB {
	PostgreSyncOnce.Do(func() {
		db, err := gorm.Open(postgres.New(postgres.Config{
			DSN:                  "host=localhost user=test password=test123 dbname=shortener-db port=5432 sslmode=disable TimeZone=Asia/Kolkata",
			PreferSimpleProtocol: true,
		}), &gorm.Config{})

		if err != nil {
			panic(err)
		}
		db.Logger = logger.Default.LogMode(logger.Info)

		sqlDB, err := db.DB()
		if err != nil {
			panic(err)
		}

		if err := sqlDB.Ping(); err != nil {
			panic(err)
		}

		if err := db.AutoMigrate(&models.ProcessedEvents{}, &models.ShortenedURL{}); err != nil {
			panic(err)
		}

		postgresClient = db
	})

	return postgresClient
}
