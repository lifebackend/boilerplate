package walleterror

import (
	"strings"
)

const (
	ForeignKeyError       = "violates foreign key constraint"
	UniqueConstraintError = "violates unique constraint"
)

func NewUnknownError(errs ...error) *AppMessage {
	if len(errs) < 1 {
		return nil
	}

	if errs[0] == nil {
		return nil
	}

	errString := toErrStr(errs)
	if strings.Contains(errString, ForeignKeyError) || strings.Contains(errString, UniqueConstraintError) {

		return NewUnexpectedForeignType(errString)
	}

	return NewInternalServiceError(errs[0], errs[1:]...)
}
