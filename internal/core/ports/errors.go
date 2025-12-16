package ports

import "fmt"

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

type MissingOrInvalidParameterError struct {
	ParameterName string
}

func (e *MissingOrInvalidParameterError) Error() string {
	return fmt.Sprintf("Missing or invalid parameter: %s", e.ParameterName)
}

type FailedOperationError struct {
	Description string
}

func (e *FailedOperationError) Error() string {
	return e.Description
}

type NotAuthorizedError struct {
	Username string
	Action   string
}

func (e *NotAuthorizedError) Error() string {
	return fmt.Sprintf("User %s is not authorized to perform action: %s", e.Username, e.Action)
}
