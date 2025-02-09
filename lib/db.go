package lib

import (
	"fmt"

	"github.com/AritroSaha10/htn25-backend-takehome/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func CreateNewDB() (*gorm.DB, error) {
	// Configure database
	var err error
	db, err = gorm.Open(sqlite.Open("htn.db"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}
	db.AutoMigrate(&model.User{}, &model.Scan{})
	return db, nil
}

func GetDB() *gorm.DB {
	return db
}
