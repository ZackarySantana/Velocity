package middleware

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/zackarysantana/velocity/internal/db"
	"golang.org/x/crypto/bcrypt"
)

// MongoDBUsernameAndPasswordUserAuthorizer is an Authorizer that uses a MongoDB
// connection to authenticate users with a username and password.
type MongoDBUsernameAndPasswordUserAuthorizer struct {
	c db.Database
}

func NewMongoDBAuthorizer(connection db.Database) MongoDBUsernameAndPasswordUserAuthorizer {
	return MongoDBUsernameAndPasswordUserAuthorizer{
		c: connection,
	}
}

func (m MongoDBUsernameAndPasswordUserAuthorizer) Auth(ctx context.Context, creds UsernameAndPasswordCredentials) (db.User, bool, error) {
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

type MongoDBAgentAuthorizer struct {
	c db.Database
}

func NewMongoDBAgentAuthorizer(connection db.Database) MongoDBAgentAuthorizer {
	return MongoDBAgentAuthorizer{
		c: connection,
	}
}

func (m MongoDBAgentAuthorizer) Auth(ctx context.Context, creds Secret) (db.Agent, bool, error) {
	var agent db.Agent
	agent, err := m.c.GetAgentBySecret(creds.Secret)
	if err != nil {
		return agent, false, fmt.Errorf("could not get entity from database: %w", err)
	}
	if agent.AgentSecret != creds.Secret {
		return agent, false, fmt.Errorf("passwords do not match: %w", err)
	}
	return agent, true, nil
}

func AuthAgentWithMongoDB(client db.Database) gin.HandlerFunc {
	return Auth[Secret, db.Agent](NewMongoDBAgentAuthorizer(client), SecretFromHeadersProvider{})
}
