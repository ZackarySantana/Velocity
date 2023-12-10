package v1types

import "github.com/zackarysantana/velocity/internal/db"

// POST /api/v1/admin/user
type PostUserRequest struct {
	Email string `json:"email"`
}
type PostUserResponse db.User

func NewPostUserRequest() interface{} {
	return &PostUserRequest{}
}
