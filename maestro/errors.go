package apihub

import "errors"

var (
	ErrLoginRequired           = errors.New("Invalid or expired token. Please log in with your ApiHub credentials.")
	ErrBadResponse             = errors.New("The response was invalid or cannot be served. For more details, execute the command with `-h`.")
	ErrTargetNotFound          = errors.New("Target not found.")
	ErrEndpointNotFound        = errors.New("You have not selected any target as default. For more details, please run `apihub target-set -h`.")
	ErrLabelExists             = errors.New("Sorry, that label has been used by another user.")
	ErrBadFormattedFile        = errors.New("Bad target data.")
	ErrCommandCancelled        = errors.New("Command Cancelled.")
	ErrFailedWritingTargetFile = errors.New("Failed trying to write the target file.")
)

type ErrorResponse struct {
	Type        string `json:"error,omitempty"`
	Description string `json:"error_description,omitempty"`
}

func (err ErrorResponse) Error() string {
	return err.Description
}

func newErrorResponse(errType, description string) ErrorResponse {
	return ErrorResponse{Type: errType, Description: description}
}

type InvalidBodyError struct {
	description error
}

func newInvalidBodyError(err error) InvalidBodyError {
	return InvalidBodyError{description: err}
}

func (err InvalidBodyError) Error() string {
	return "Request body could not be parsed: " + err.description.Error()
}

type InvalidHostError struct {
	description error
}

func newInvalidHostError(err error) InvalidHostError {
	return InvalidHostError{description: err}
}

func (err InvalidHostError) Error() string {
	return "You either have not selected any target or it is invalid: " + err.description.Error()
}

type RequestError struct {
	description error
}

func newRequestError(err error) RequestError {
	return RequestError{description: err}
}

func (err RequestError) Error() string {
	return "Failed to connect to ApiHub server: " + err.description.Error()
}

type ResponseError struct {
	description error
}

func newResponseError(err error) ResponseError {
	return ResponseError{description: err}
}

func (err ResponseError) Error() string {
	return err.description.Error()
}

type UnauthorizedError struct {
	description error
}

func newUnauthorizedError(err error) UnauthorizedError {
	return UnauthorizedError{description: err}
}

func (err UnauthorizedError) Error() string {
	return err.description.Error()
}
