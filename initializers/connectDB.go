package initializers

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB(config *Config) {
	var err error
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", config.DBUserName, config.DBUserPassword, config.DBName)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the Database")
	}
	fmt.Println("Connected Successfully to the Database")
}
