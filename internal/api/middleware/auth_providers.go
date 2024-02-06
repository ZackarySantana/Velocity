package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// MultiCredentialProvider combines multiple CredentialProviders into one.
// It attempts them one by one until one succeeds.
type MultiCredentialProvider[T any] struct {
	providers []CredentialProvider[T]
}

// NewMultiCredentialProvider creates a new MultiCredentialProvider with the given
// providers.
func NewMultiCredentialProvider[T any](providers ...CredentialProvider[T]) CredentialProvider[T] {
	return MultiCredentialProvider[T]{
		providers: providers,
	}
}

func (m MultiCredentialProvider[T]) Get(ctx *gin.Context) (T, error, error) {
	var creds T

	for _, provider := range m.providers {
		creds, invalidErr, err := provider.Get(ctx)
		if err != nil || invalidErr != nil {
			// TODO: Surface these errors
			continue
		}
		return creds, nil, nil
	}

	return creds, fmt.Errorf("no providers could parse your credentials"), nil
}

// CreateUsernameAndPasswordMultiProvider creates a MultiCredentialProvider for
// all available UsernameAndPasswordCredentials providers.
func CreateUsernameAndPasswordMultiProvider() CredentialProvider[UsernameAndPasswordCredentials] {
	return NewMultiCredentialProvider[UsernameAndPasswordCredentials](
		UsernameAndPasswordFromJSONBodyProvider{},
		UsernameAndPasswordFromHeadersProvider{},
	)
}

type UsernameAndPasswordFromJSONBodyProvider struct{}

func (u UsernameAndPasswordFromJSONBodyProvider) Get(ctx *gin.Context) (UsernameAndPasswordCredentials, error, error) {
	var creds UsernameAndPasswordCredentials
	err := ctx.ShouldBindJSON(creds)
	if err != nil {
		return UsernameAndPasswordCredentials{}, nil, fmt.Errorf("could not bind json: %w", err)
	}
	if creds.Username == "" || creds.Password == "" {
		return UsernameAndPasswordCredentials{}, fmt.Errorf("username or password is empty"), nil
	}
	return creds, nil, nil
}

type UsernameAndPasswordFromHeadersProvider struct{}

func (u UsernameAndPasswordFromHeadersProvider) Get(ctx *gin.Context) (UsernameAndPasswordCredentials, error, error) {
	username := ctx.GetHeader("username")
	password := ctx.GetHeader("password")
	if username == "" || password == "" {
		return UsernameAndPasswordCredentials{}, fmt.Errorf("username or password is empty"), nil
	}
	return UsernameAndPasswordCredentials{
		Username: username,
		Password: password,
	}, nil, nil
}

type SecretFromHeadersProvider struct{}

func (s SecretFromHeadersProvider) Get(ctx *gin.Context) (Secret, error, error) {
	secret := ctx.GetHeader("secret")
	if secret == "" {
		return Secret{}, fmt.Errorf("secret is empty"), nil
	}
	return Secret{
		Secret: secret,
	}, nil, nil
}
