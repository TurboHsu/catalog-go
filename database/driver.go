package database

import (
	"fmt"

	"catalog-go/config"
	"catalog-go/database/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() (err error) {
	switch config.CONFIG.Database.Type {
	case "sqlite3":
		DB, err = gorm.Open(sqlite.Open(config.CONFIG.Database.Path), &gorm.Config{})
	default:
		err = fmt.Errorf("unsupported database type: %s", config.CONFIG.Database.Type)
	}
	return
}

func MigrateDatabase() (err error) {
	for _, model := range model.ALL {
		err = DB.AutoMigrate(&model)
		if err != nil {
			return
		}
	}
	return
}

func CloseDatabase() (err error) {
	sqlDB, err := DB.DB()
	if err != nil {
		return
	}
	err = sqlDB.Close()
	return
}
