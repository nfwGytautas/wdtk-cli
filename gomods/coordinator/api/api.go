package api

import (
	"log"

	"github.com/nfwGytautas/mstk/gomods/common-api"
)

// Database
var dbConn common.DatabaseConnection

// ========================================================================
// PUBLIC
// ========================================================================

/*
Sets up the api package
*/
func Setup() {
	log.Println("Preparing API package")

	// Open connection
	dbConn = common.DatabaseConnection{}
	dbConn.Initialize(common.DatabaseConnectionConfig{
		DCS:             dcs,
		MigrateCallback: nil,
	})

	log.Println("API package ready")
}

// ========================================================================
// PRIVATE
// ========================================================================

const dcs = "mstk:mstk123@tcp(coordinator-db:3306)/coordinator_db?charset=utf8mb4&parseTime=True&loc=Local"
