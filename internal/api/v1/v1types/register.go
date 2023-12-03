package v1types

import "github.com/zackarysantana/velocity/internal/db"

type RegisterUserRequest struct {
	Email string `json:"email"`
}

type RegisterUserResponse db.User
