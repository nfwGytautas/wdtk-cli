package auth

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// ========================================================================
// PUBLIC
// ========================================================================

/*
Middleware for authenticating

Usage:
r.Use(auth.AuthenticationMiddleware())
*/
func AuthenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		info, err := parseToken(c)
		if err != nil {
			log.Println(err)
			c.String(http.StatusInternalServerError, "Token")
			c.Abort()
			return
		}

		if !info.valid {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		c.Next()
	}
}

/*
Middleware for authorization

Usage:
r.Use(auth.AuthorizationMiddleware([]string{"role"}))
*/
func AuthorizationMiddleware(roles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		info, err := parseToken(c)
		if err != nil {
			log.Println(err)
			c.String(http.StatusInternalServerError, "Token")
			c.Abort()
		}

		valid := false
		for _, role := range roles {
			if role == info.role {
				valid = true
				break
			}
		}

		if !valid {
			c.String(http.StatusUnauthorized, "Access denied")
			c.Abort()
			return
		}

		c.Next()
	}
}

// ========================================================================
// PRIVATE
// ========================================================================

/*
Struct for containing token info
*/
type tokenInfo struct {
	valid bool
	id    uint
	role  string
}

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

/*
Parse a token from gin context
*/
func parseToken(c *gin.Context) (tokenInfo, error) {
	result := tokenInfo{}
	result.valid = false

	tokenString := c.Query("token")

	if tokenString == "" {
		// Token empty check if it is inside Authorization header
		tokenString = c.Request.Header.Get("Authorization")

		// Since this is bearer token we need to parse the token out
		if len(strings.Split(tokenString, " ")) == 2 {
			tokenString = strings.Split(tokenString, " ")[1]
		} else {
			return result, errors.New("invalid request")
		}
	}

	jwtToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.Secret), nil
	})

	if err != nil {
		return result, err
	}

	// Token valid fill token information
	claims, ok := jwtToken.Claims.(jwt.MapClaims)

	if !ok || !jwtToken.Valid {
		return result, nil
	}

	// User id
	uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["user_id"]), 10, 32)
	if err != nil {
		return result, nil
	}

	result.id = uint(uid)

	// Role
	result.role = claims["role"].(string)

	result.valid = true
	return result, nil
}
