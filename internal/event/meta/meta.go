package meta

import (
	"strconv"

	"github.com/zackarysantana/velocity/internal/db"
)

func CreateApplyIndexes(user db.User) map[string]string {
	return map[string]string{
		"user_id": user.Id.Hex(),
		"type":    "system_wide",
	}
}

func CreateUser(creator db.User, created db.User) map[string]string {
	return map[string]string{
		"creator_id": creator.Id.Hex(),
		"created_id": created.Id.Hex(),
		"super_user": strconv.FormatBool(created.UserPermission.SuperUser),
		"type":       "user",
	}
}
