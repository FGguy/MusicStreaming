package data

import "fmt"

type UserNotFoundError struct {
	username string
}

func (e *UserNotFoundError) Error() string {
	return fmt.Sprintf("No user with username: %s", e.username)
}
