package auth

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/nfwGytautas/mstk/gomods/common-api"
	"golang.org/x/crypto/bcrypt"
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

const TokenLifespan = 60 // Lifespan in minutes
const APISecret = "MSTK_API_SECRET_TEST"
const DBConnectionString = "mstk:mstk123@tcp(auth-db:3306)/auth_db?charset=utf8mb4&parseTime=True&loc=Local"

var dbConn common.DatabaseConnection

/*
Setup authentication database connection
*/
func Setup() {
	dbConn = common.DatabaseConnection{}
	dbConn.Initialize(common.DatabaseConnectionConfig{
		DCS: DBConnectionString,
		MigrateCallback: func(d *gorm.DB) {
			d.AutoMigrate(&User{})
		},
	})
}

/*
Adds GIN handlers for authentication
*/
func AddRoutes(r *gin.Engine) {
	v := r.Group("/auth", common.RequireDatabaseConnectionMiddleware(&dbConn))

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

func meHandler(c *gin.Context) {
	var u User

	info, err := common.ParseToken(c)
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

func verifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
