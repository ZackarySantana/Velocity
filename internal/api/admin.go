package api

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zackarysantana/velocity/internal/api/middleware"
	"github.com/zackarysantana/velocity/internal/db"
	"github.com/zackarysantana/velocity/internal/event"
	"github.com/zackarysantana/velocity/internal/event/meta"
	"golang.org/x/crypto/bcrypt"
)

func (a *Api) AddAdminRoutes() {
	admin := a.Group("/admin")
	admin.Use(middleware.AuthUsernameAndPasswordUserWithMongoDB(a.db), middleware.OnlySuperUsers)

	admin.POST("/user/create", a.CreateUser)
	admin.GET("/indexes/apply", a.ApplyIndexes)
}

type CreateUserRequest struct {
	User db.User `json:"user"`
}

func (c *CreateUserRequest) Validate() error {
	if c.User.Username == "" {
		return errors.New("username is required")
	}
	if len(c.User.Username) < 8 || len(c.User.Username) > 24 {
		return errors.New("username must between 8 and 24 characters")
	}
	if c.User.Password == "" {
		return errors.New("password is required")
	}
	if len(c.User.Password) < 8 || len(c.User.Password) > 24 {
		return errors.New("password must between 8 and 24 characters")
	}
	if c.User.Email == "" {
		return errors.New("email is required")
	}
	if len(c.User.Email) < 3 {
		return errors.New("email is too short")
	}
	atIndex := strings.Index(c.User.Email, "@")
	if atIndex <= 0 {
		return errors.New("email is invalid and needs to include an @")
	}
	dotIndex := strings.LastIndex(c.User.Email, ".")
	if dotIndex <= atIndex+1 || dotIndex == len(c.User.Email)-1 {
		return errors.New("email is invalid and needs to include a . after @")
	}
	return nil
}

// POST /admin/user/create
func (a *Api) CreateUser(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(&gin.Error{
			Err:  fmt.Errorf("could not parse your request: %v", err.Error()),
			Type: gin.ErrorTypePublic,
		})
		ctx.Abort()
		return
	}

	if err := req.Validate(); err != nil {
		ctx.Error(&gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic,
		})
		ctx.Abort()
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte(req.User.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.Error(&gin.Error{
			Err:  fmt.Errorf("error hashing password: %v", err.Error()),
			Type: gin.ErrorTypePublic,
		})
		ctx.Abort()
		return
	}
	req.User.Password = string(password)

	user, err := a.db.CreateUser(ctx, req.User)
	if err != nil {
		ctx.Error(&gin.Error{
			Err:  err,
			Type: gin.ErrorTypePrivate,
		})
		ctx.Error(&gin.Error{
			Err:  fmt.Errorf("error creating user: %v", err.Error()),
			Type: gin.ErrorTypePublic,
		})
		ctx.Abort()
		return
	}

	ctx.JSON(200, gin.H{
		"user_id": user.Id.Hex(),
		"message": "user created",
	})
}

// GET /admin/indexes/apply
func (a *Api) ApplyIndexes(ctx *gin.Context) {
	err := a.db.ApplyIndexes(ctx)
	if err != nil {
		ctx.Error(&gin.Error{
			Err:  err,
			Type: gin.ErrorTypePrivate,
		})
		ctx.Error(&gin.Error{
			Err:  fmt.Errorf("error applying indexes: %v", err.Error()),
			Type: gin.ErrorTypePublic,
		})
		ctx.Abort()
		return
	}

	user := middleware.MustGetAuthArtifact[db.User](ctx)
	err = a.es.SendEvent(ctx, event.Event{
		EventType: event.EventTypeIndexesApplied,
		Metadata:  meta.CreateApplyIndexes(user),
	})
	if err != nil {
		ctx.Error(&gin.Error{
			Err:  err,
			Type: gin.ErrorTypePrivate,
		})
		fmt.Println("Testing", err)
	}

	ctx.JSON(200, gin.H{
		"message": "indexes applied",
	})
}
