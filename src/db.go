package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func GetDb() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("data/exporter.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect sqlite database")
	}
	return db
}
