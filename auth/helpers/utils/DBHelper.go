package utils

import (
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var postgresClient *gorm.DB = nil
var PostgreSyncOnce sync.Once

func ConnectToPostgres() *gorm.DB {
	PostgreSyncOnce.Do(func() {
		db, err := gorm.Open(postgres.New(postgres.Config{
			DSN:                  "host=localhost user=test password=test123 dbname=url-shortner port=5432 sslmode=disable TimeZone=Asia/Kolkata",
			PreferSimpleProtocol: true,
		}), &gorm.Config{})

		if err != nil {
			panic(err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			panic(err)
		}

		if err := sqlDB.Ping(); err != nil {
			panic(err)
		}

		postgresClient = db
	})

	return postgresClient
}
