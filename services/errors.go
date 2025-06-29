package services

import (
	"errors"
	"fmt"
)

var (
	ErrTaskNotFound      = errors.New("task not found")
	ErrUnauthorizedAccess = errors.New("unauthorized access to task")
	ErrInvalidInput      = errors.New("invalid input data")
	ErrTaskAlreadyCompleted = errors.New("task is already completed")
	ErrDueDateInPast     = errors.New("due date cannot be in the past")
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

func (v ValidationErrors) Error() string {
	if len(v.Errors) == 0 {
		return "validation failed"
	}
	return fmt.Sprintf("validation failed: %s", v.Errors[0].Message)
}

func NewValidationError(field, message string) ValidationErrors {
	return ValidationErrors{
		Errors: []ValidationError{
			{Field: field, Message: message},
		},
	}
}

func (v *ValidationErrors) AddError(field, message string) {
	v.Errors = append(v.Errors, ValidationError{
		Field:   field,
		Message: message,
	})
}

func (v ValidationErrors) HasErrors() bool {
	return len(v.Errors) > 0
}