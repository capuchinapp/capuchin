package util

import (
	"fmt"
)

func ErrTrace(msg string, e error) error {
	return fmt.Errorf("%s: %w", msg, e)
}
