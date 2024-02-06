package middleware

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
)

// Authorizer has two generic types, T which is the type of credentials, and E
// which is the type of the artifact that is returned from the Auth method.
// Authorizers can be found in auth_authorizers.go
type Authorizer[T any, E any] interface {
	// Auth attempts to authenticate the given credentials 'T' and returns
	// an artifact 'E' if successful. If the credentials are invalid, it
	// returns false. If an error occurs, it returns an error.
	Auth(context.Context, T) (E, bool, error)
}

// CredentialProvider has one generic type, T which is the type of credentials.
// CredentialProviders can be found in auth_providers.go
type CredentialProvider[T any] interface {
	// Get attempts to get the credentials from the given context. If the
	// credentials are invalid, it returns the first error. If an error
	// 	occurs, it returns the second error.
	Get(r *gin.Context) (T, error, error)
}

// Auth creates a middleware function that uses the provided Authorizer and
// CredentialProvider to authenticate requests. It handles errors, cancel the
// if needed, and sets the artifact returned from the Authorizer in the context
// as "auth_artifact".
func Auth[T any, E any](authorizer Authorizer[T, E], provider CredentialProvider[T]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		creds, invalidErr, err := provider.Get(ctx)
		if err != nil {
			ctx.Error(&gin.Error{
				Err:  err,
				Type: gin.ErrorTypePrivate,
			})
			ctx.Error(&gin.Error{
				Err:  errors.New("could not parse your credentials"),
				Type: gin.ErrorTypePublic,
			})
			ctx.Abort()
			return
		}
		if invalidErr != nil {
			ctx.Error(&gin.Error{
				Err:  errors.New("your credentials are invalid"),
				Type: gin.ErrorTypePublic,
			})
			ctx.Abort()
			return
		}
		artifact, authed, err := authorizer.Auth(ctx, creds)
		if err != nil {
			ctx.Error(&gin.Error{
				Err:  err,
				Type: gin.ErrorTypePrivate,
			})
			ctx.Error(&gin.Error{
				Err:  errors.New("there was an error authenticating you"),
				Type: gin.ErrorTypePublic,
			})
			ctx.Abort()
			return
		}
		if !authed {
			ctx.Error(&gin.Error{
				Err:  errors.New("you are not authorized"),
				Type: gin.ErrorTypePublic,
			})
			ctx.Abort()
			return
		}
		ctx.Set("auth_artifact", artifact)
		ctx.Next()
	}
}

type UsernameAndPasswordCredentials struct {
	Username string
	Password string
}
