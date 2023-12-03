package v1types

import "github.com/zackarysantana/velocity/internal/db"

type CreateUserRequest struct {
	Email string `json:"email"`
}

type CreateUserResponse db.User
