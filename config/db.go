package config

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/jinzhu/gorm"
)

var gormConn *gorm.DB

// InitDB ...
func InitDB() {
	if gormConn != nil && gormConn.DB() != nil && gormConn.DB().Ping() == nil {
		return
	}

	// Try to connect to the database
	conn, err := gorm.Open(os.Getenv("DB_DIALECT"), os.Getenv("DB_CONNECTION"))
	if err != nil {
		logrus.Fatal("Could not connect to the database")
	}

	// Setting connection pool
	conn.DB().SetMaxIdleConns(100)
	conn.DB().SetMaxOpenConns(10000)

	// Store the connection in package variable for furher request
	gormConn = conn
}

// GetDatabaseConnection returns gorm connection
func GetDatabaseConnection() *gorm.DB {
	// Check if a connection allready exists
	return gormConn
}
