package bootstrap

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/MicaelaJofre/gocourse_domain/domain"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func DBConecction() (*gorm.DB, error) {
	dbPortStr := os.Getenv("DATABASE_PORT")
	if dbPortStr == "" {
		log.Fatal("DATABASE_PORT environment variable not set.")
	}
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		log.Fatalf("Invalid DATABASE_PORT value: %v", err)
	}

	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
	os.Getenv("DATABASE_ENROLLMENT"),
	os.Getenv("DATABASE_PASSWORD"),
	os.Getenv("DATABASE_HOST"),
	dbPort,
	os.Getenv("DATABASE_NAME"))

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if os.Getenv("DATABASE_DEBUG") == "true" {
	db = db.Debug()
	}

	if os.Getenv("DATABASE_MIGRATE") == "true" {
		if err := db.AutoMigrate(&domain.Enrollment{}); err != nil {
			return nil, err
		}
	}
	return db, nil
}

func InitLogger() *log.Logger {
	l := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	return l
}