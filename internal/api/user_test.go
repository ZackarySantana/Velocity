package api_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zackarysantana/velocity/internal/api"
	"github.com/zackarysantana/velocity/internal/cli/logger"
	"github.com/zackarysantana/velocity/internal/db"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateUserRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     api.CreateUserRequest
		wantErr string
	}{
		{
			name: "valid request",
			req: api.CreateUserRequest{
				User: db.User{
					Username: "testtest",
					Password: "password",
					Email:    "test@test.com",
				},
			},
			wantErr: "",
		},
		{
			name: "no username",
			req: api.CreateUserRequest{
				User: db.User{
					Username: "",
					Password: "password",
					Email:    "test@test.com",
				},
			},
			wantErr: "username is required",
		},
		{
			name: "too short username",
			req: api.CreateUserRequest{
				User: db.User{
					Username: "test",
					Password: "password",
					Email:    "test@test.com",
				},
			},
			wantErr: "username must between 8 and 24 characters",
		},
		{
			name: "too long username",
			req: api.CreateUserRequest{
				User: db.User{
					Username: "testtesttesttesttesttesttest",
					Password: "password",
					Email:    "test@test.com",
				},
			},
			wantErr: "username must between 8 and 24 characters",
		},
		{
			name: "no password",
			req: api.CreateUserRequest{
				User: db.User{
					Username: "testtest",
					Password: "",
					Email:    "test@test.com",
				},
			},
			wantErr: "password is required",
		},
		{
			name: "too short password",
			req: api.CreateUserRequest{
				User: db.User{
					Username: "testtest",
					Password: "pass",
					Email:    "test@test.com",
				},
			},
			wantErr: "password must between 8 and 24 characters",
		},
		{
			name: "too long password",
			req: api.CreateUserRequest{
				User: db.User{
					Username: "testtest",
					Password: "passwordpasswordpasswordpasswordpassword",
					Email:    "test@test.com",
				},
			},
			wantErr: "password must between 8 and 24 characters",
		},
		{
			name: "too short email",
			req: api.CreateUserRequest{
				User: db.User{
					Username: "testtest",
					Password: "password",
					Email:    "kc",
				},
			},
			wantErr: "email is too short",
		},
		{
			name: "no @ in email",
			req: api.CreateUserRequest{
				User: db.User{
					Username: "testtest",
					Password: "password",
					Email:    "test.com",
				},
			},
			wantErr: "email is invalid and needs to include an @",
		},
		{
			name: "no . in email",
			req: api.CreateUserRequest{
				User: db.User{
					Username: "testtest",
					Password: "password",
					Email:    "test@testcom",
				},
			},
			wantErr: "email is invalid and needs to include a . after @",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()

			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	d, err := db.NewMockWithUsers()
	require.NoError(t, err)
	a := api.CreateApi(logger.NewCollectLogger(), d)
	a.AddUserRoutes()

	tests := []struct {
		name       string
		admin      bool
		req        *http.Request
		body       api.CreateUserRequest
		resCode    int
		resMessage string
		resError   string
	}{
		{
			name:  "create regular user",
			admin: true,
			body: api.CreateUserRequest{
				User: db.User{
					Username: "testtest",
					Password: "password",
					Email:    "new_user@test.com",
				},
			},
			resCode:    200,
			resMessage: "user created",
		},
		{
			name:  "create regular user with no permissions",
			admin: false,
			body: api.CreateUserRequest{
				User: db.User{
					Username: "testtest",
					Password: "password",
					Email:    "new_user@test.com",
				},
			},
			resCode:  401,
			resError: "user is not a super user",
		},
		{
			name:  "create super user",
			admin: true,
			body: api.CreateUserRequest{
				User: db.User{
					Username: "testtest",
					Password: "password",
					Email:    "new_user@test.com",
					UserPermission: db.UserPermission{
						SuperUser: true,
					},
				},
			},
			resCode:    200,
			resMessage: "user created",
		},
		{
			name:  "create super user with no permissions",
			admin: false,
			body: api.CreateUserRequest{
				User: db.User{
					Username: "testtest",
					Password: "password",
					Email:    "new_user@test.com",
					UserPermission: db.UserPermission{
						SuperUser: true,
					},
				},
			},
			resCode:  401,
			resError: "user is not a super user",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := NewRequest("POST", "/user/create", IsAdminRequest(tt.admin), WithJSONBody(tt.body))
			require.NoError(t, err)

			a.ServeHTTP(w, req)

			assert.Equal(t, tt.resCode, w.Code)

			body, err := io.ReadAll(w.Body)
			require.NoError(t, err)
			var res map[string]string
			err = json.Unmarshal(body, &res)
			require.NoError(t, err)

			assert.Equal(t, tt.resMessage, res["message"])
			assert.Equal(t, tt.resError, res["error"])

			// If the response code is 200, then the user_id should not be empty.
			if tt.resCode == 200 {
				assert.NotEmpty(t, res["user_id"])
			} else {
				assert.Empty(t, res["user_id"])
			}

			// If the response is 200, it should be added to our database.
			if tt.resCode == 200 {
				user, err := d.GetUserByUsername(context.Background(), tt.body.User.Username)
				require.NoError(t, err)
				assert.Equal(t, tt.body.User.Username, user.Username)
				assert.Equal(t, tt.body.User.Email, user.Email)
				assert.Equal(t, tt.body.User.UserPermission.SuperUser, user.UserPermission.SuperUser)
				assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(tt.body.User.Password)))
			}
		})
	}
}
