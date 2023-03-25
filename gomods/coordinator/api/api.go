package api

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Database
var db *gorm.DB

// ========================================================================
// PUBLIC
// ========================================================================

/*
Sets up the api package
*/
func Setup(dcs string) {
	var err error

	log.Println("Preparing API package")

	// Open connection
	db, err = gorm.Open(mysql.Open(dcs), &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}

	// Setup table migration
	db.AutoMigrate(&Service{}, &Endpoint{}, &Shard{})

	log.Println("API package ready")
}
