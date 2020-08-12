package main

import (
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

//SetupDB Pega a conecção com o banco de daods
func SetupDB() (*gorm.DB, error) {
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	var db *gorm.DB

	settings := "host=" + host + " user=" + user + " password=" + pass + " dbname=" + dbname + " sslmode=disable"
	db, err := gorm.Open("postgres", settings)

	if err != nil {
		return nil, err
	}

	err = db.DB().Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
