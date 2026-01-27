package ports

import "fmt"

// NotFoundError indicates that a requested resource was not found.
// It maps to Subsonic API error code 70.
type NotFoundError struct {
	Message string
}

// Error implements the error interface for NotFoundError.
func (e *NotFoundError) Error() string {
	return e.Message
}

// MissingOrInvalidParameterError indicates that a required parameter is missing or invalid.
// It maps to Subsonic API error code 10.
type MissingOrInvalidParameterError struct {
	ParameterName string
}

// Error implements the error interface for MissingOrInvalidParameterError.
func (e *MissingOrInvalidParameterError) Error() string {
	return fmt.Sprintf("Missing or invalid parameter: %s", e.ParameterName)
}

// FailedOperationError indicates that an operation failed for unspecified reasons.
// It maps to Subsonic API error code 0 (generic error).
type FailedOperationError struct {
	Description string
}

// Error implements the error interface for FailedOperationError.
func (e *FailedOperationError) Error() string {
	return e.Description
}

// NotAuthorizedError indicates that the user is not authorized to perform the requested action.
// It maps to Subsonic API error code 50.
type NotAuthorizedError struct {
	Username string
	Action   string
}

// Error implements the error interface for NotAuthorizedError.
func (e *NotAuthorizedError) Error() string {
	return fmt.Sprintf("User %s is not authorized to perform action: %s", e.Username, e.Action)
}
