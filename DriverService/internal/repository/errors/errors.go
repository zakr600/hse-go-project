package errors

import "fmt"

type NotFoundError struct {
	Key string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("not found: %s", e.Key)
}

type DuplicateKeyError struct {
	Key string
}

func (e DuplicateKeyError) Error() string {
	return fmt.Sprintf("duplicate key: %s", e.Key)
}
