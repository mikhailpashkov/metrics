package models

import "errors"

var (
	ErrMetricsIsNil        = errors.New("metrics is nil")
	ErrNotFound            = errors.New("not found")
	ErrConflict            = errors.New("conflict")
	ErrInvalidInput        = errors.New("invalid input")
	ErrDeleteAllNotAllowed = errors.New("delete all not allowed")
	ErrInternal            = errors.New("internal error")
)
