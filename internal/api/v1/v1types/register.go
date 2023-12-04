package v1types

import "github.com/zackarysantana/velocity/internal/db"

type PostRegisterUserRequest struct {
	Email string `json:"email"`
}

type PostRegisterUserResponse db.User
