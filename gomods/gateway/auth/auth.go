package auth

import (
	"fmt"
	"log"
	"net/http"
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

/*
Struct containing information about authentication configuration
*/
type AuthConfig struct {
	Name     string
	Password string
	URL      string
	DBName   string

	TokenLifespan int    // Minutes that a JWT token lasts
	Secret        string // API secret
}

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

// Keep config information global
var config AuthConfig

// Database
var db *gorm.DB

/*
Setup authentication database connection
*/
func Setup(cfg AuthConfig) {
	var err error

	log.Println("Trying to connect to auth database")
	config = cfg
	common.APISecret = config.Secret

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.Name, config.Password, config.URL, config.DBName)

	log.Printf("Connection string: %s", dsn)

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}

	db.AutoMigrate(&User{})

	log.Println("Database UP and running")
}

/*
Adds GIN handlers for authentication
*/
func AddRoutes(r *gin.Engine) {
	v := r.Group("/auth")

	v.POST("/login", loginHandler)
	v.POST("/register", registerHandler)

	vP := v.Group("/", common.AuthenticationMiddleware())
	vP.GET("/me", meHandler)
}

// ========================================================================
// PRIVATE
// ========================================================================

/*
Generate a an access token for the specified user id
*/
func generateToken(user *User) (string, error) {
	claims := jwt.MapClaims{}

	claims["authorized"] = true
	claims["user_id"] = user.ID
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(config.TokenLifespan)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.Secret))

}

func loginHandler(c *gin.Context) {
	// Request model
	input := struct {
		Identifier string
		Password   string
	}{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := User{}

	// Get username
	err := db.Model(User{}).Where("identifier = ?", input.Identifier).Take(&u).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify password
	err = verifyPassword(input.Password, u.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
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
		Identifier string
		Password   string
	}{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := User{}
	u.Identifier = input.Identifier
	u.Password = input.Password
	u.Role = "new"

	err := db.Create(&u).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	err = db.First(&u, info.ID).Error
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
