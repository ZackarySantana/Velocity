package db

import (
	"context"

	"github.com/zackarysantana/velocity/internal/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	Id     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	APIKey string             `bson:"api_key" json:"api_key"`

	Email string `bson:"email" json:"email"`
}

type Permissions struct {
	Id     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	APIKey string             `bson:"api_key" json:"api_key"`

	UserId primitive.ObjectID `bson:"user_id" json:"user_id"`

	Admin bool `bson:"admin" json:"admin"`
}

func (c *Connection) CreateUser(ctx context.Context, email string) (*User, error) {
	apiKey, err := api.GenerateAPIKey()
	if err != nil {
		return nil, err
	}

	user := User{
		APIKey: apiKey,
		Email:  email,
	}

	r, err := c.col("users").InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	user.Id = r.InsertedID.(primitive.ObjectID)

	return &user, nil
}

func (c *Connection) CreateAdminUser(ctx context.Context, email string) (*User, error) {
	apiKey, err := api.GenerateAPIKey()
	if err != nil {
		return nil, err
	}

	user := User{
		APIKey: apiKey,
		Email:  email,
	}

	err = c.UseSessionWithOptions(ctx, nil, func(ctx mongo.SessionContext) error {
		if err := ctx.StartTransaction(); err != nil {
			return err
		}

		r, err := c.col("users").InsertOne(ctx, user)
		if err != nil {
			_ = ctx.AbortTransaction(context.Background())
			return err
		}
		user.Id = r.InsertedID.(primitive.ObjectID)

		_, err = c.col("permissions").InsertOne(ctx, Permissions{
			APIKey: apiKey,
			UserId: user.Id,
			Admin:  true,
		})
		if err != nil {
			_ = ctx.AbortTransaction(context.Background())
			return err
		}

		return ctx.CommitTransaction(context.Background())
	})
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *Connection) GetUser(ctx context.Context, query interface{}) (*User, error) {
	var user User
	return &user, c.col("users").FindOne(ctx, query).Decode(&user)
}

func (c *Connection) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return c.GetUser(ctx, bson.M{"email": email})
}

func (c *Connection) GetUserByAPIKey(ctx context.Context, apiKey string) (*User, error) {
	return c.GetUser(ctx, bson.M{"api_key": apiKey})
}

func (c *Connection) GetPermissions(ctx context.Context, query interface{}) (*Permissions, error) {
	var permissions Permissions
	return &permissions, c.col("permissions").FindOne(ctx, query).Decode(&permissions)
}

func (c *Connection) GetPermissionsByAPIKey(ctx context.Context, apiKey string) (*Permissions, error) {
	return c.GetPermissions(ctx, bson.M{"api_key": apiKey})
}

func (c *Connection) HasAdminUsers(ctx context.Context) (bool, error) {
	count, err := c.col("permissions").CountDocuments(ctx, bson.M{"admin": true})
	return count > 0, err
}
