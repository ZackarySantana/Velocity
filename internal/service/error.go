package service

import (
	"errors"
	"strings"
)

func ParseError(err error) error {
	if err == nil {
		return nil
	}

	switch strings.Trim(err.Error(), "\n") {
	case ErrEmptyQueue.Error():
		return ErrEmptyQueue
	case ErrInvalidId.Error():
		return ErrInvalidId
	}

	return err
}

func ParseErrorMsg(err string) error {
	return ParseError(errors.New(err))
}
