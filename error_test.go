package estree

import (
	"errors"
)

// hasError indicates one of errs satisfies errors.Is for expect.
func hasError(expect error, errs ...error) bool {
	for _, err := range errs {
		if errors.Is(err, expect) {
			return true
		}
	}
	return false
}
