package types

import (
	"fmt"
	"github.com/go-errors/errors"
)

func NewError(errType, msg string) error {
	message := fmt.Sprintf("%s:%s", errType, msg)
	return errors.New(errors.Errorf(message))
}

func ErrorStackTrace(err error) string {
	return fmt.Sprintf(err.(*errors.Error).ErrorStack())
}
