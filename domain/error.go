package domain

import errors "github.com/rotisserie/eris"

var (
	ErrGraphCycle      = errors.New("graph cycle detected")
	ErrRecordNotFound  = errors.New("record not found")
	ErrNotImplemented  = errors.New("not implemented")
	ErrDuplicateRecord = errors.New("duplicate record")
	ErrBodyAttribute   = errors.New("body attribute error")
)
