package config

import (
	"github.com/samber/oops"
	"github.com/zackarysantana/velocity/src/catcher"
)

type Validator interface {
	validateSyntax() error
	validateIntegrity(*Config) error

	error(int) oops.OopsErrorBuilder
}

func Validate(v Validator, c *Config) error {
	catcher := catcher.New()
	catcher.Wrap(v.validateSyntax(), "validating syntax")
	catcher.Wrap(v.validateIntegrity(c), "validating integrity")
	return catcher.Resolve()
}

func ValidateSyntax(v Validator) error {
	return v.error(0).Wrapf(v.validateSyntax(), "validating syntax")
}

func ValidateIntegrity(v Validator, c *Config) error {
	return v.error(0).Wrapf(v.validateIntegrity(c), "validating integrity")
}
