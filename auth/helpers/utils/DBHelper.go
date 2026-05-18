package utils

import (
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var postgresClient *gorm.DB = nil
var PostgreSyncOnce sync.Once

func ConnectToPostgres() *gorm.DB {
	// https://github.com/go-gorm/postgres
	if postgresClient != nil {
		return postgresClient
	}

	PostgreSyncOnce.Do(
		func() {
			db, err := gorm.Open(postgres.New(postgres.Config{
				DSN:                  "user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai",
				PreferSimpleProtocol: true, // disables implicit prepared statement usage
			}), &gorm.Config{})

			if err != nil {
				panic(err)
			}

			postgresClient = db

		},
	)

	return postgresClient

}
