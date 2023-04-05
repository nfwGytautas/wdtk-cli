package database

import (
	"github.com/nfwGytautas/mstk/lib/gdev/database"
	"gorm.io/gorm"
)

// PUBLIC TYPES
// ========================================================================

// PRIVATE TYPES
// ========================================================================

var dbConn database.DatabaseConnection

/*
Connection string
*/
const dcs = "mstk:mstk123@tcp(locator-db:3306)/locator_db?charset=utf8mb4&parseTime=True&loc=Local"

// PUBLIC FUNCTIONS
// ========================================================================

/*
Setup the database
*/
func Setup() {
	dbConn = database.DatabaseConnection{}
	dbConn.Initialize(database.DatabaseConnectionConfig{
		DCS: dcs,
		MigrateCallback: func(d *gorm.DB) {
			d.AutoMigrate(&Service{}, &Endpoint{}, &Role{})
		},
	})
}

/*
Get a service with the specified name
*/
func GetService(service string) (Service, error) {
	var out Service
	result := dbConn.DB.Preload("Endpoints").Where("name = ?", service).First(&out)
	return out, result.Error
}

/*
Create a service
*/
func CreateService(service *Service) error {
	err := dbConn.DB.Create(service).Error
	return err
}

// PRIVATE FUNCTIONS
// ========================================================================
