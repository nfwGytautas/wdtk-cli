package database

import (
	"github.com/nfwGytautas/mstk/gomods/api/common-api"
	"gorm.io/gorm"
)

// ========================================================================
// PUBLIC
// ========================================================================

/*
Setup the database
*/
func Setup() {
	dbConn = common.DatabaseConnection{}
	dbConn.Initialize(common.DatabaseConnectionConfig{
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

// ========================================================================
// PRIVATE
// ========================================================================

var dbConn common.DatabaseConnection

/*
Connection string
*/
const dcs = "mstk:mstk123@tcp(locator-db:3306)/locator_db?charset=utf8mb4&parseTime=True&loc=Local"
