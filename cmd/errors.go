package cmd

import "errors"

var ErrNotFound = errors.New("record not found")

func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}
