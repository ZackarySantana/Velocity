package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type Mock struct {
	Users map[string]User
}

func NewMock() Database {
	return &Mock{
		Users: make(map[string]User),
	}
}

func NewMockWithUsers() (Database, error) {
	adminPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	userPassword, err := bcrypt.GenerateFromPassword([]byte("user"), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &Mock{
		Users: map[string]User{
			"admin": {
				Id:       primitive.NewObjectID(),
				Username: "admin",
				Password: string(adminPassword),
				Email:    "admin@test.com",
				UserPermission: UserPermission{
					SuperUser: true,
				},
			},
			"user": {
				Id:       primitive.NewObjectID(),
				Username: "user",
				Password: string(userPassword),
				Email:    "user@test.com",
				UserPermission: UserPermission{
					SuperUser: false,
				},
			},
		},
	}, nil
}

func (m *Mock) GetUserByUsername(ctx context.Context, username string) (User, error) {
	if user, ok := m.Users[username]; ok {
		return user, nil
	}
	return User{}, ErrNoEntity
}

func (m *Mock) CreateUser(ctx context.Context, user User) (User, error) {
	user.Id = primitive.NewObjectID()
	m.Users[user.Username] = user
	return user, nil
}

func (m *Mock) GetAgentBySecret(ctx context.Context, agentSecret string) (Agent, error) {
	return Agent{}, nil
}
