package database

import "gorm.io/gorm"

/*
This file contains tables that are defined in the locator database
*/

// ========================================================================
// PUBLIC
// ========================================================================

/*
GORM service struct
*/
type Service struct {
	gorm.Model

	Name      string     `gorm:"unique" json:"Name"`
	URL       string     `json:"URL"`
	Endpoints []Endpoint `json:"Endpoints"`
}

/*
GORM endpoint struct
*/
type Endpoint struct {
	gorm.Model

	ServiceID uint `json:"-"`

	Name         string `json:"Name"`
	Method       string `json:"Method"`
	AllowedRoles []Role `json:"AllowedRoles"`
}

/*
GORM role struct
*/
type Role struct {
	gorm.Model

	Name        string `json:"Name"`
	Description string `json:"Description"`
}
