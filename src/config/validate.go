package config

import "github.com/samber/oops"

type Validator interface {
	validateSyntax() error
	validateIntegrity(*Config) error
}

func validate(v Validator) error {
	oops := oops.With("object", v)
	if err := v.validateSyntax(); err != nil {
		return oops.Wrapf(err, "validating syntax")
	}
	if err := v.validateIntegrity(nil); err != nil {
		return oops.Wrapf(err, "validating integrity")
	}
	return nil
}
