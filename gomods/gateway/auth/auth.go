package auth

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/nfwGytautas/mstk/gomods/common-api"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ========================================================================
// PUBLIC
// ========================================================================

//TODO: Add pepper and salt to authentication

/*
Struct describing the User table
*/
type User struct {
	gorm.Model
	Identifier string // Identifier for users e.g. email, username, etc.
	Password   string // Salt&Pepper hashed password
	Role       string // Role of the user (for applications that don't use Authorization this is useless)
}

// Database
var dbConnection struct {
	mx sync.RWMutex
	db *gorm.DB
}

const TokenLifespan = 60 // Lifespan in minutes
const APISecret = "MSTK_API_SECRET_TEST"
const DBConnectionString = "mstk:mstk123@tcp(auth-db:3306)/auth_db?charset=utf8mb4&parseTime=True&loc=Local"

/*
Setup authentication database connection
*/
func Setup() {
	dbConnection.db = nil
	go checkDBConnection()
}

/*
Adds GIN handlers for authentication
*/
func AddRoutes(r *gin.Engine) {
	v := r.Group("/auth", dbMiddleware())

	v.POST("/login", loginHandler)
	v.POST("/register", registerHandler)

	vP := v.Group("/", common.AuthenticationMiddleware())
	vP.GET("/me", meHandler)
}

// ========================================================================
// PRIVATE
// ========================================================================

/*
Middleware for making sure database is online
*/
func dbMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		dbConnection.mx.Lock()

		if dbConnection.db == nil {
			log.Println("Database not ready yet returning 503")
			c.Status(http.StatusServiceUnavailable)
			c.Abort()
			dbConnection.mx.Unlock()
			return
		}

		dbConnection.mx.Unlock()
		c.Next()
	}
}

/*
Continually perform checks on the database connection
*/
func checkDBConnection() {
	var err error
	log.Println("DB heartbeat started")

	for range time.Tick(time.Second * 5) {
		if dbConnection.db == nil {
			dbConnection.mx.Lock()

			log.Println("Trying to connect to auth database")

			dbConnection.db, err = gorm.Open(mysql.Open(DBConnectionString), &gorm.Config{})
			if err != nil {
				dbConnection.db = nil
				dbConnection.mx.Unlock()
				continue
			}

			dbConnection.db.AutoMigrate(&User{})

			log.Println("Database UP and running")

			dbConnection.mx.Unlock()
		}
	}
}

/*
Generate a an access token for the specified user id
*/
func generateToken(user *User) (string, error) {
	claims := jwt.MapClaims{}

	claims["authorized"] = true
	claims["user_id"] = user.ID
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(TokenLifespan)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(APISecret))

}

func loginHandler(c *gin.Context) {
	// Request model
	input := struct {
		Identifier string `json:"identifier" binding:"required"`
		Password   string `json:"password" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := User{}

	// Get username
	dbConnection.mx.Lock()
	err := dbConnection.db.Model(User{}).Where("identifier = ?", input.Identifier).Take(&u).Error
	dbConnection.mx.Unlock()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify password
	err = verifyPassword(input.Password, u.Password)
	if err != nil || err == bcrypt.ErrMismatchedHashAndPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Credentials correct, create token and return it
	token, err := generateToken(&u)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func registerHandler(c *gin.Context) {
	// Request model
	input := struct {
		Identifier string `json:"identifier" binding:"required"`
		Password   string `json:"password" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := User{}
	u.Identifier = input.Identifier

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u.Password = string(hash)
	u.Role = "new"

	dbConnection.mx.Lock()
	err = dbConnection.db.Create(&u).Error
	dbConnection.mx.Unlock()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "registration success"})
}

func meHandler(c *gin.Context) {
	var u User

	info, err := common.ParseToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dbConnection.mx.Lock()
	err = dbConnection.db.First(&u, info.ID).Error
	dbConnection.mx.Unlock()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Remove password fields
	u.Password = ""

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": u})
}

func verifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
