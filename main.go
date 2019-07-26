package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"webserver/config"
	"webserver/models"
	app "webserver/server"
)

func main() {
	// Loads the env variables
	err := godotenv.Load()
	if err != nil {
		logrus.Warning("Cannot load .env file!")
	}
	config.InitDB()

	// Check database connection
	db := config.GetDatabaseConnection()
	defer db.Close()

	// Prints the version and the address of our api to the console
	logrus.Info("Version is ", os.Getenv("PORT"))
	logrus.Info("Starting Server on http://localhost:", os.Getenv("PORT"))

	// Set log level
	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		logrus.SetLevel(logrus.ErrorLevel)
	case "info":
		logrus.SetLevel(logrus.ErrorLevel)
	case "warn":
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.SetLevel(logrus.ErrorLevel)
	}

	// Creates the database schema
	migrateDatabase()

	// Server router on given port and attach the cors headers
	server := app.NewServer()
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), server))
}

func migrateDatabase() {
	db := config.GetDatabaseConnection()

	// Migrate the given tables
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Confession{})
}
