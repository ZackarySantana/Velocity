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
	APIKey string             `bson:"api_key,omitempty" json:"api_key"`

	Email string `bson:"email,omitempty" json:"email"`
}

type Permissions struct {
	Id     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	APIKey string             `bson:"api_key,omitempty" json:"api_key"`

	UserId primitive.ObjectID `bson:"user_id,omitempty" json:"user_id"`

	Admin bool `bson:"admin,omitempty" json:"admin"`
}

func (c *Connection) ApplyUserIndexes(ctx context.Context) error {
	_, err := c.col("users").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{"email": 1, "api_key": 1},
	})
	return err
}

func (c *Connection) ApplyPermissionIndexes(ctx context.Context) error {
	_, err := c.col("permissions").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{"user_id": 1, "api_key": 1},
	})
	return err
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

func (c *Connection) InsertUser(ctx context.Context, email string) (*User, error) {
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

func (c *Connection) InsertAdminUser(ctx context.Context, email string) (*User, error) {
	apiKey, err := api.GenerateAPIKey()
	if err != nil {
		return nil, err
	}

	var user *User

	err = c.UseSessionWithOptions(ctx, nil, func(ctx mongo.SessionContext) error {
		if err := ctx.StartTransaction(); err != nil {
			return err
		}

		user, err = c.InsertUser(ctx, email)
		if err != nil {
			_ = ctx.AbortTransaction(context.Background())
			return err
		}

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

	return user, nil
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
