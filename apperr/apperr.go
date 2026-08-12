package apperr

import (
	"errors"
	"fmt"
)

type Kind string

const (
	KindNotFound     Kind = "not_found"
	KindConflict     Kind = "conflict"
	KindValidation   Kind = "validation"
	KindUnauthorized Kind = "unauthorized"
	KindInternal     Kind = "internal"
)

type Error struct {
	Kind    Kind
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return string(e.Kind)
}

func (e *Error) Unwrap() error {
	return e.Err
}

func Conflict(msg string) error {
	return &Error{Kind: KindConflict, Message: msg}
}

func NotFound(msg string) error {
	return &Error{Kind: KindNotFound, Message: msg}
}

func Validation(msg string) error {
	return &Error{Kind: KindValidation, Message: msg}
}

func Unauthorized(msg string) error {
	return &Error{Kind: KindUnauthorized, Message: msg}
}

func Internal(msg string, err error) error {
	if msg == "" && err != nil {
		msg = err.Error()
	}
	return &Error{Kind: KindInternal, Message: msg, Err: err}
}

func IsKind(err error, kind Kind) bool {
	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr.Kind == kind
	}
	return false
}

// WrapInternal returns err unchanged if it is already an *Error; otherwise wraps it as Internal.
func WrapInternal(msg string, err error) error {
	if err == nil {
		return nil
	}
	var appErr *Error
	if errors.As(err, &appErr) {
		return err
	}
	if msg == "" {
		return Internal(err.Error(), err)
	}
	return Internal(fmt.Sprintf("%s: %v", msg, err), err)
}
