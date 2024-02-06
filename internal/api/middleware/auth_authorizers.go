package middleware

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/zackarysantana/velocity/internal/db"
	"golang.org/x/crypto/bcrypt"
)

// MongoDBUsernameAndPasswordAuthorizer is an Authorizer that uses a MongoDB
// connection to authenticate users with a username and password.
type MongoDBUsernameAndPasswordAuthorizer struct {
	c db.Database
}

func NewMongoDBAuthorizer(connection db.Database) MongoDBUsernameAndPasswordAuthorizer {
	return MongoDBUsernameAndPasswordAuthorizer{
		c: connection,
	}
}

func (m MongoDBUsernameAndPasswordAuthorizer) Auth(ctx context.Context, creds UsernameAndPasswordCredentials) (db.User, bool, error) {
	var user db.User
	user, err := m.c.GetUserByUsername(creds.Username)
	if err != nil {
		return user, false, fmt.Errorf("could not get entity from database: %w", err)
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		return user, false, fmt.Errorf("passwords do not match: %w", err)
	}
	return user, true, nil
}

// AuthUsernameAndPasswordUserWithMongoDB creates a middleware function that
// authenticates requests with a username and password using a MongoDB connection.
// The providers is uses are all that are available for UsernameAndPasswordCredentials.
func AuthUsernameAndPasswordUserWithMongoDB(client db.Database) gin.HandlerFunc {
	return Auth[UsernameAndPasswordCredentials, db.User](NewMongoDBAuthorizer(client), CreateUsernameAndPasswordMultiProvider())
}
