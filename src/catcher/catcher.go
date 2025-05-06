package catcher

import (
	"errors"
	"fmt"

	"github.com/samber/oops"
)

type Catcher struct {
	errs []error
}

func New() *Catcher {
	return &Catcher{}
}

func (c *Catcher) Catch(err error) {
	if err == nil {
		return
	}
	c.errs = append(c.errs, err)
}

func (c *Catcher) Wrap(err error, msg string, a ...any) {
	if err == nil {
		return
	}
	c.Catch(joinOopsErrorsChain(fmt.Errorf(msg, a...), err))
}

func (c *Catcher) New(msg string, a ...any) {
	if msg == "" {
		return
	}
	c.Catch(fmt.Errorf(msg, a...))
}

func (c *Catcher) When(cond bool, msg string, a ...any) {
	if !cond {
		return
	}
	c.New(msg, a...)
}

func (c *Catcher) Resolve() error {
	if len(c.errs) == 0 {
		return nil
	}
	var builder oops.OopsErrorBuilder
	for _, e := range c.errs {
		oopsErr, ok := oops.AsOops(e)
		if !ok {
			continue
		}
		for k, v := range oopsErr.Context() {
			builder = builder.With(k, v)
		}
	}
	return builder.Wrap(errors.Join(c.errs...))
}

func joinOopsErrorsChain(errs ...error) error {
	return oopsBuilderWithAllContext(errs...).Wrap(Join(errs...))
}

func oopsBuilderWithAllContext(errs ...error) oops.OopsErrorBuilder {
	var builder oops.OopsErrorBuilder
	for _, e := range errs {
		oopsErr, ok := oops.AsOops(e)
		if !ok {
			continue
		}
		for k, v := range oopsErr.Context() {
			builder = builder.With(k, v)
		}
	}
	return builder
}
