package middleware

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/zackarysantana/velocity/internal/db"
	"golang.org/x/crypto/bcrypt"
)

// DatabaseUsernameAndPasswordUserAuthorizer is an Authorizer that uses a Database
// interface to authenticate users with a username and password.
type DatabaseUsernameAndPasswordUserAuthorizer struct {
	d db.Database
}

func NewDatabaseAuthorizer(connection db.Database) DatabaseUsernameAndPasswordUserAuthorizer {
	return DatabaseUsernameAndPasswordUserAuthorizer{
		d: connection,
	}
}

func (m DatabaseUsernameAndPasswordUserAuthorizer) Auth(ctx context.Context, creds UsernameAndPasswordCredentials) (db.User, bool, error) {
	var user db.User
	user, err := m.d.GetUserByUsername(ctx, creds.Username)
	if err != nil {
		if err == db.ErrNoEntity {
			return user, false, nil
		}
		return user, false, fmt.Errorf("could not get user from database: %w", err)
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		return user, false, nil
	}
	return user, true, nil
}

// curl command with username and password in body
// curl -X GET http://localhost:8080/user/ping -d '{"username":"test","password":"password"}' -H "Content-Type: application/json"

// AuthUsernameAndPasswordUserWithMongoDB creates a middleware function that
// authenticates requests with a username and password using a Database interface.
// The providers is uses are all that are available for UsernameAndPasswordCredentials.
func AuthUsernameAndPasswordUserWithMongoDB(client db.Database) gin.HandlerFunc {
	return Auth[UsernameAndPasswordCredentials, db.User](NewDatabaseAuthorizer(client), UsernameAndPasswordFromHeadersProvider{})
}

type DatabaseAgentAuthorizer struct {
	d db.Database
}

func NewDatabaseAgentAuthorizer(connection db.Database) DatabaseAgentAuthorizer {
	return DatabaseAgentAuthorizer{
		d: connection,
	}
}

func (m DatabaseAgentAuthorizer) Auth(ctx context.Context, creds Secret) (db.Agent, bool, error) {
	var agent db.Agent
	agent, err := m.d.GetAgentBySecret(ctx, creds.Secret)
	if err != nil {
		if err == db.ErrNoEntity {
			return agent, false, nil
		}
		return agent, false, fmt.Errorf("could not get agent from database: %w", err)
	}
	if agent.AgentSecret != creds.Secret {
		return agent, false, fmt.Errorf("passwords do not match: %w", err)
	}
	return agent, true, nil
}

func AuthAgentWithMongoDB(client db.Database) gin.HandlerFunc {
	return Auth[Secret, db.Agent](NewDatabaseAgentAuthorizer(client), SecretFromHeadersProvider{})
}
