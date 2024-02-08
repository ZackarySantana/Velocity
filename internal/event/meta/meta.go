package meta

import "github.com/zackarysantana/velocity/internal/db"

func CreateApplyIndexes(user db.User) map[string]string {
	return map[string]string{
		"user_id": user.Id.Hex(),
		"type":    "system_wide",
	}
}
