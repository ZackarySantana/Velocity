package v1types

import "github.com/zackarysantana/velocity/internal/db"

// POST /api/v1/first_time_register
type PostFirstTimeRegisterRequest struct {
	Email string `json:"email"`
}
type PostFirstTimeRegisterResponse db.User
