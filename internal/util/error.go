package util

import "fmt"

func Wrap(errp error, format string, args ...interface{}) {
	if errp != nil {
		s := fmt.Sprintf(format, args...)
		errp = fmt.Errorf("%s: %w", s, errp)
	}
}
