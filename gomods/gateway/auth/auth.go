package auth

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/nfwGytautas/mstk/lib/gdev/database"
	jwtu "github.com/nfwGytautas/mstk/lib/gdev/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// PUBLIC TYPES
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

// PRIVATE TYPES
// ========================================================================

var tokenLifespan = 60
var externalSecret = ""

const dbBConnectionString = "mstk:mstk123@tcp(auth-db:3306)/auth_db?charset=utf8mb4&parseTime=True&loc=Local"

var dbConn database.DatabaseConnection

// PUBLIC FUNCTIONS
// ========================================================================

/*
Setup authentication database connection
*/
func Setup() {
	dbConn = database.DatabaseConnection{}
	dbConn.Initialize(database.DatabaseConnectionConfig{
		DCS: dbBConnectionString,
		MigrateCallback: func(d *gorm.DB) {
			d.AutoMigrate(&User{})

			// Delete MSTK user
			var user User
			var err error

			result := d.Where("role = ?", "_mstk").Delete(&user)

			if result.Error != nil {
				log.Panic(result.Error)
			}

			// Create MSTK user
			user.Identifier = os.Getenv("MSTK_USER")
			user.Password, err = createPasswordHash(os.Getenv("MSTK_PSW"))
			user.Role = "_mstk"

			if err != nil {
				log.Panic(err)
			}

			if user.Identifier == "" || user.Password == "" {
				panic(50)
			}

			d.Create(&user)
		},
	})

	// Assign values
	var err error
	tokenLifespan, err = strconv.Atoi(os.Getenv("TOKEN_LIFESPAN"))
	externalSecret = os.Getenv("API_SECRET")

	if err != nil {
		log.Panic(err)
	}
}

/*
Adds GIN handlers for authentication
*/
func AddRoutes(r *gin.Engine) {
	v := r.Group("/auth", database.RequireDatabaseConnectionMiddleware(&dbConn))

	v.POST("/login", loginHandler)
	v.POST("/register", registerHandler)

	vP := v.Group("/", jwtu.AuthenticationMiddleware())
	vP.GET("/me", meHandler)
}

// PRIVATE FUNCTIONS
// ========================================================================

/*
Generate a an access token for the specified user id
*/
func generateToken(user *User) (string, error) {
	claims := jwt.MapClaims{}

	claims["authorized"] = true
	claims["user_id"] = user.ID
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(tokenLifespan)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(externalSecret))

}

/*
Handler for handling login requests
*/
func loginHandler(c *gin.Context) {
	// Request model
	input := struct {
		Identifier string `json:"identifier" binding:"required"`
		Password   string `json:"password" binding:"required"`
	}{}

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get username
	u := User{}

	err = dbConn.DB.Model(User{}).Where("identifier = ?", input.Identifier).Take(&u).Error
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

/*
Handler for register requests
*/
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

	hash, err := createPasswordHash(input.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u.Password = hash
	u.Role = "new"

	err = dbConn.DB.Model(User{}).Where("identifier = ?", input.Identifier).Take(&u).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = dbConn.DB.Create(&u).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "registration success"})
}

/*
Handler that returns the current logged in user details
*/
func meHandler(c *gin.Context) {
	var u User

	info, err := jwtu.ParseToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = dbConn.DB.First(&u, info.ID).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Remove password fields
	u.Password = ""

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": u})
}

/*
Verify that the passed password string matches the hashed password
*/
func verifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

/*
Creates a hash from password
*/
func createPasswordHash(psw string) (string, error) {
	// TODO: Salt & pepper

	hash, err := bcrypt.GenerateFromPassword([]byte(psw), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}
