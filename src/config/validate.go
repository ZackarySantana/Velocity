package config

import (
	"github.com/samber/oops"
)

type Validator interface {
	validateSyntax() error
	validateIntegrity(*Config) error

	Error() oops.OopsErrorBuilder
}

func validate(v Validator, c *Config) error {
	if err := validateSyntax(v); err != nil {
		return err
	}
	if err := validateIntegrity(v, c); err != nil {
		return err
	}
	return nil
}

func validateSyntax(v Validator) error {
	if err := v.validateSyntax(); err != nil {
		return v.Error().Wrapf(err, "validating syntax")
	}
	return nil
}

func validateIntegrity(v Validator, c *Config) error {
	if err := v.validateIntegrity(c); err != nil {
		return v.Error().Wrapf(err, "validating integrity")
	}
	return nil
}
