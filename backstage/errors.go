package main

import "errors"

var (
	ErrLoginRequired           = errors.New("Invalid or expired token. Please log in with your Backstage credentials.")
	ErrBadRequest              = errors.New("The request was invalid or cannot be served.")
	ErrLabelExists             = errors.New("Sorry, that label has been used by another user.")
	ErrLabelNotFound           = errors.New("Sorry, that label does not exist.")
	ErrBadFormattedFile        = errors.New("Bad target data.")
	ErrCommandCancelled        = errors.New("Command Cancelled.")
	ErrFailedWritingTargetFile = errors.New("Failed trying to write the target file.")
	ErrFailedConnectingServer  = errors.New("Failed to connect to the server. Please check if the target is correct.")
	ErrEndpointNotFound        = errors.New("You have not selected any target as default. For more details, please run `backstage target-set -h`.")
)

// The HTTPError type is a http representation of error.
type HTTPError struct {
	ErrorDescription string `json:"error_description"`
	Url              string `json:"url"`
}

func (err *HTTPError) Error() string {
	return err.ErrorDescription
}
