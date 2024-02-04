package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Authorizer[T any] interface {
	Auth(context.Context, T) (bool, error)
}

type AuthProvider[T any] interface {
	Get(r *http.Request) (T, error, error)
}

type UsernameAndPassword struct {
	Username string
	Password string
}

type MongoDBUsernameAndPasswordAuthorizer struct {
	c          mongo.Client
	database   string
	collection string
}

func NewMongoDBAuthorizer(connection mongo.Client, database, collection string) MongoDBUsernameAndPasswordAuthorizer {
	return MongoDBUsernameAndPasswordAuthorizer{
		c:          connection,
		database:   database,
		collection: collection,
	}
}

func (m MongoDBUsernameAndPasswordAuthorizer) Auth(ctx context.Context, creds UsernameAndPassword) (bool, error) {
	fmt.Println(creds)
	fmt.Println("here")
	test, err := m.c.Database(m.database).Collection(m.collection).Find(ctx, bson.M{"username": creds.Username})
	if err != nil {
		return false, fmt.Errorf("could not get entity from database: %w", err)
	}
	if test == nil {
		return false, nil
	}
	type entity struct {
		Password string `bson:"password"`
	}
	var e entity
	if err = test.Decode(&e); err != nil {
		return false, fmt.Errorf("could not decode entity from database: %w", err)
	}
	if err = bcrypt.CompareHashAndPassword([]byte(e.Password), []byte(creds.Password)); err != nil {
		return false, fmt.Errorf("passwords do not match: %w", err)
	}
	return true, nil
}

type UsernameAndPasswordFromBodyProvider struct{}

func (u UsernameAndPasswordFromBodyProvider) Get(r *http.Request) (UsernameAndPassword, error, error) {
	fmt.Println("there")
	return UsernameAndPassword{}, nil, nil
}

func Auth[T any](authorizer Authorizer[T], provider AuthProvider[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		creds, invalidErr, err := provider.Get(c.Request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "an error occurred"})
			c.Abort()
			return
		}
		if invalidErr != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "your credentials are invalid"})
			c.Abort()
			return
		}
		authed, err := authorizer.Auth(c, creds)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "an error occurred"})
			c.Abort()
			return
		}
		if !authed {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "you are not authorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func AuthWithMongoDBAndUsernameAndPasswordFromBody(client mongo.Client, database, collection string) gin.HandlerFunc {
	return Auth[UsernameAndPassword](NewMongoDBAuthorizer(client, database, collection), UsernameAndPasswordFromBodyProvider{})
}
