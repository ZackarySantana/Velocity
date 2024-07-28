package catcher

import (
	"errors"
	"fmt"
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
	c.Catch(errors.Join(err, fmt.Errorf(msg, a...)))
}

func (c *Catcher) Error(msg string, a ...any) {
	if msg == "" {
		return
	}
	c.Catch(fmt.Errorf(msg, a...))
}

func (c *Catcher) Resolve() error {
	if len(c.errs) == 0 {
		return nil
	}
	return errors.Join(c.errs...)
}
