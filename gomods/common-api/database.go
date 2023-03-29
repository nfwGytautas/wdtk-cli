package common

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ========================================================================
// PUBLIC
// ========================================================================

/*
Migration callback function
*/
type MigrationFn func(*gorm.DB)

/*
Configuration for database connection
*/
type DatabaseConnectionConfig struct {
	DCS             string      // Connection string
	MigrateCallback MigrationFn // Callback function for migrating database content
}

/*
Database connection struct that provides automatic health checks
*/
type DatabaseConnection struct {
	mx       sync.RWMutex
	DB       *gorm.DB
	cfg      DatabaseConnectionConfig
	migrated bool
}

/*
Initialize a database connection with a given database connection string, returns
error on error, nil otherwise
*/
func (dc *DatabaseConnection) Initialize(cfg DatabaseConnectionConfig) error {
	dc.DB = nil
	dc.cfg = cfg
	dc.migrated = false

	dc.connect()

	return nil
}

/*
Middleware for gin that requires a database connection
*/
func RequireDatabaseConnectionMiddleware(dc *DatabaseConnection) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Avoid processing requests since service unavailable
		if !dc.ping() {
			ctx.Abort()
			ctx.Status(http.StatusServiceUnavailable)
			return
		}

		if !dc.migrated && dc.cfg.MigrateCallback != nil {
			// Migrate
			dc.mx.Lock()
			dc.cfg.MigrateCallback(dc.DB)
			dc.migrated = true
			dc.mx.Unlock()
		}

		ctx.Next()
	}
}

// ========================================================================
// PRIVATE
// ========================================================================

/*
Connect to database
*/
func (dc *DatabaseConnection) connect() {
	var err error

	dc.DB, err = gorm.Open(mysql.Open(dc.cfg.DCS), &gorm.Config{})
	if err != nil {
		// Failed to open connection
		log.Println(err)
	}
}

/*
Ping the database connection to check if we are still online
*/
func (dc *DatabaseConnection) ping() bool {
	db, err := dc.DB.DB()

	if err != nil {
		// We failed to get a DB instance?
		log.Println(err)
		return false
	}

	err = db.Ping()
	return err == nil
}
