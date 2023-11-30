package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zackarysantana/velocity/internal/db"
)

var (
	authenticators = []func(*gin.Context) (*db.User, *authError){
		authBySessionToken,
	}
)

type authError struct {
	Msg  string
	Code int
}

// Auth is a simple authentication middleware with an option to
// check for if the user is an admin or not
// Commented code is not for gin, covert it to gin code
func Auth(c *gin.Context) {
	for _, authenticator := range authenticators {
		user, err := authenticator(c)
		if err != nil {
			c.AbortWithStatusJSON(err.Code, gin.H{"error": err.Msg})
			return
		}

		if user != nil {
			c.Set("user", user)
			c.Next()
			return
		}
	}

	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized no authenticator found"})
}

func authBySessionToken(c *gin.Context) (*db.User, *authError) {
	db := GetDB(c)

	sessionToken, err := c.Cookie("session_token")
	if err != nil {
		return nil, nil
	}

	user, err := db.GetUserBySessionToken(sessionToken)
	if err != nil {
		return nil, &authError{"Unauthorized session_token", http.StatusUnauthorized}
	}

	expireTime, timeParseErr := time.Parse(time.RFC3339, user.SessionExpires)
	if timeParseErr != nil {
		return nil, &authError{"Unauthorized session_token", http.StatusUnauthorized}
	}

	if time.Now().After(expireTime) {
		return nil, &authError{"Unauthorized session_token", http.StatusUnauthorized}
	}

	// get user from db

	return nil, nil
}
