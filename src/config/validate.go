package config

import (
	"fmt"

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

func ValidateSyntaxMany(v []Validator) error {
	if v == nil {
		return nil
	}
	catcher := catcher.New()
	for i, val := range v {
		catcher.Catch(val.error(i).Wrapf(val.validateSyntax(), fmt.Sprintf("[index=%d]", i)))
	}
	return catcher.Resolve()
}

func ValidateIntegrityMany(v []Validator, c *Config) error {
	if v == nil {
		return nil
	}
	catcher := catcher.New()
	for i, val := range v {
		catcher.Catch(val.error(i).Wrapf(val.validateIntegrity(c), fmt.Sprintf("[index=%d]", i)))
	}
	return catcher.Resolve()
}

func ValidateIntegrity(v Validator, c *Config) error {
	return v.error(0).Wrapf(v.validateIntegrity(c), "validating integrity")
}

func toValidators(v interface{}) []Validator {
	if v == nil {
		return nil
	}
	var validators []Validator

	switch v := v.(type) {
	case *TestSection:
		for _, test := range *v {
			validators = append(validators, &test)
		}
	case *ImageSection:
		for _, image := range *v {
			validators = append(validators, &image)
		}
	case *JobSection:
		for _, job := range *v {
			validators = append(validators, &job)
		}
	case *RoutineSection:
		for _, routine := range *v {
			validators = append(validators, &routine)
		}
	default:
		panic(fmt.Sprintf("unsupported type: %T", v))
	}

	return validators
}
