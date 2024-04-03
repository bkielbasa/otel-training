package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

type TemperatureRecord struct {
	gorm.Model
	Address      string
	TemperatureC float64
	TemperatureF float64
}

type WindInfo struct {
	gorm.Model              // Adds fields ID, CreatedAt, UpdatedAt, DeletedAt
	WindSpeed     float64   `gorm:"column:wind_speed"`
	WindDirection string    `gorm:"column:wind_direction"`
	Location      string    `gorm:"column:location"`
	LocalTime     time.Time `gorm:"column:localtime"`
}

var db *gorm.DB
var err error

const defaultPort = "8080"

func initDB() {
	var (
		host     = os.Getenv("DB_HOST")
		port     = os.Getenv("DB_PORT")
		user     = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASSWORD")
		dbname   = os.Getenv("DB_NAME")
	)

	dbURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", host, port, user, dbname, password)
	db, err = gorm.Open("postgres", dbURL)
	if err != nil {
		panic("failed to connect to database " + dbURL)
	}

	db = db.AutoMigrate(&TemperatureRecord{}).AutoMigrate(&WindInfo{})
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	initDB()
	defer db.Close()

	http.HandleFunc("/temperature/{address}", temperatureHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func temperatureHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handling")

	tempC := 22.0
	tempF := tempC*9/5 + 32

	record := TemperatureRecord{Address: r.PathValue("address"), TemperatureC: tempC, TemperatureF: tempF}
	db.Create(&record)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(record)
}
