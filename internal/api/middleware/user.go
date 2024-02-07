package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/zackarysantana/velocity/internal/db"
)

func OnlySuperUsers(ctx *gin.Context) {
	user := MustGetAuthArtifact[db.User](ctx)
	if !user.UserPermission.SuperUser {
		ctx.Error(&gin.Error{
			Err:  errors.New("user is not a super user"),
			Type: gin.ErrorTypePublic,
			Meta: 401,
		})
		ctx.Abort()
		return
	}
	ctx.Next()
}
