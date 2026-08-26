package signalhouse

import "fmt"

// ValidationError is returned when a required parameter is missing.
type ValidationError struct {
	Field string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("missing required parameter: %s", e.Field)
}

// APIError represents a non-2xx response from the SignalHouse API.
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("signalhouse api error (status %d): %s", e.StatusCode, e.Message)
}
