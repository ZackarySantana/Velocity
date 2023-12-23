package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zackarysantana/velocity/internal/db"
)

var (
	authenticators = []func(*gin.Context, db.Connection) (*db.User, *authError){
		func(ctx *gin.Context, db db.Connection) (*db.User, *authError) {
			return authByAPIKey(ctx, db.GetUserByAPIKey)
		},
	}

	adminAuthenticators = []func(*gin.Context, db.Connection) (*db.Permissions, *authError){
		func(ctx *gin.Context, db db.Connection) (*db.Permissions, *authError) {
			return authByAPIKey(ctx, db.GetPermissionsByAPIKey)
		},
	}
)

type authError struct {
	Msg  string
	Code int
}

func Auth(db db.Connection) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, authenticator := range authenticators {
			user, err := authenticator(c, db)
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

		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized- no authentication found"})
	}
}

func GetUser(c *gin.Context) *db.User {
	user, _ := c.Get("user")
	return user.(*db.User)
}

func AdminAuth(client db.Connection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var permissions *db.Permissions
		var err *authError
		for _, authenticator := range adminAuthenticators {
			permissions, err = authenticator(c, client)
			if err != nil {
				c.AbortWithStatusJSON(err.Code, gin.H{"error": err.Msg})
				return
			}

			if permissions != nil {
				break
			}
		}

		if permissions != nil && permissions.Admin {
			c.Set("permissions", permissions)
			c.Next()
			return
		}

		if permissions != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized- you are not an admin"})
			return
		}

		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized- no authentication found"})
	}
}

func GetPermissions(c *gin.Context) *db.Permissions {
	permissions, _ := c.Get("permissions")
	return permissions.(*db.Permissions)
}

func authByAPIKey[T any](c *gin.Context, retrieve func(ctx context.Context, apiKey string) (*T, error)) (*T, *authError) {
	apiKey := c.Request.Header.Get("Authorization")
	if apiKey == "" {
		return nil, nil
	}

	s := strings.Split(apiKey, "Bearer ")
	if len(s) != 2 {
		return nil, &authError{"Invalid authorization header", http.StatusBadRequest}
	}
	apiKey = s[1]

	resource, err := retrieve(c, apiKey)
	if err != nil || resource == nil {
		return nil, &authError{"Unauthorized api key", http.StatusUnauthorized}
	}

	return resource, nil
}
