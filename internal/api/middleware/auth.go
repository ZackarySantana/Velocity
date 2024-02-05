package middleware

import (
	"context"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Authorizer[T any] interface {
	Auth(context.Context, T) (bool, error)
}

type AuthProvider[T any] interface {
	Get(r *gin.Context) (T, error, error)
}

type UsernameAndPassword struct {
	Username string
	Password string
}

type MongoDBUsernameAndPasswordAuthorizer struct {
	c          *mongo.Client
	database   string
	collection string
}

func NewMongoDBAuthorizer(connection *mongo.Client, database, collection string) MongoDBUsernameAndPasswordAuthorizer {
	return MongoDBUsernameAndPasswordAuthorizer{
		c:          connection,
		database:   database,
		collection: collection,
	}
}

func (m MongoDBUsernameAndPasswordAuthorizer) Auth(ctx context.Context, creds UsernameAndPassword) (bool, error) {
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

type UsernameAndPasswordFromJSONBodyProvider struct{}

func (u UsernameAndPasswordFromJSONBodyProvider) Get(ctx *gin.Context) (UsernameAndPassword, error, error) {
	var creds UsernameAndPassword
	err := ctx.ShouldBindJSON(creds)
	if err != nil {
		return UsernameAndPassword{}, nil, fmt.Errorf("could not bind json: %w", err)
	}
	if creds.Username == "" || creds.Password == "" {
		return UsernameAndPassword{}, fmt.Errorf("username or password is empty"), nil
	}
	return creds, nil, nil
}

func Auth[T any](authorizer Authorizer[T], provider AuthProvider[T]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		creds, invalidErr, err := provider.Get(ctx)
		if err != nil {
			ctx.Error(&gin.Error{
				Err:  err,
				Type: gin.ErrorTypePrivate,
			})
			ctx.Error(&gin.Error{
				Err:  errors.New("could not parse your credentials"),
				Type: gin.ErrorTypePublic,
			})
			ctx.Abort()
			return
		}
		if invalidErr != nil {
			ctx.Error(&gin.Error{
				Err:  errors.New("your credentials are invalid"),
				Type: gin.ErrorTypePublic,
			})
			ctx.Abort()
			return
		}
		authed, err := authorizer.Auth(ctx, creds)
		if err != nil {
			ctx.Error(&gin.Error{
				Err:  err,
				Type: gin.ErrorTypePrivate,
			})
			ctx.Error(&gin.Error{
				Err:  errors.New("there was an error authenticating you"),
				Type: gin.ErrorTypePublic,
			})
			ctx.Abort()
			return
		}
		if !authed {
			ctx.Error(&gin.Error{
				Err:  errors.New("you are not authorized"),
				Type: gin.ErrorTypePublic,
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func AuthWithMongoDBAndUsernameAndPasswordFromJSONBody(client *mongo.Client, database, collection string) gin.HandlerFunc {
	return Auth[UsernameAndPassword](NewMongoDBAuthorizer(client, database, collection), UsernameAndPasswordFromJSONBodyProvider{})
}
