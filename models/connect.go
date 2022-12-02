package models

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDataBase() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	Dbdriver := os.Getenv("DB_DRIVER")
	DbHost := os.Getenv("DB_HOST")
	DbUser := os.Getenv("DB_USER")
	DbPassword := os.Getenv("DB_PASSWORD")
	DbName := os.Getenv("DB_NAME")
	DbPort := os.Getenv("DB_PORT")
	DBTimeZone := os.Getenv("DB_TIME_ZONE")

	var URLDatabase string = ""
	// connect from mysql
	if Dbdriver == "mysql" {
		URLDatabase = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
		DB, err = gorm.Open(mysql.Open(URLDatabase), &gorm.Config{})
	} else if Dbdriver == "postgres" {
		URLDatabase = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s", DbHost, DbUser, DbPassword, DbName, DbPort, DBTimeZone)
		DB, err = gorm.Open(postgres.Open(URLDatabase), &gorm.Config{})
	}

	if err != nil {
		fmt.Println(URLDatabase)
		fmt.Println("Cannot connect to database ", Dbdriver)
	} else {
		fmt.Println("connected to the database ", Dbdriver)
	}

	DB.AutoMigrate(&Customer{})
	DB.AutoMigrate(&Transaction{})
}
