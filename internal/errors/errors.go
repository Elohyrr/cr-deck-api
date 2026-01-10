package errors

import "fmt"

// APIError represents errors from external APIs
type APIError struct {
	StatusCode int
	Message    string
	Endpoint   string
	RetryAfter int
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error [%d] on %s: %s", e.StatusCode, e.Endpoint, e.Message)
}

func (e *APIError) IsNotFound() bool    { return e.StatusCode == 404 }
func (e *APIError) IsRateLimited() bool { return e.StatusCode == 429 }
func (e *APIError) IsServerError() bool { return e.StatusCode >= 500 }

// DBError represents database errors
type DBError struct {
	Operation string
	Table     string
	Err       error
}

func (e *DBError) Error() string {
	return fmt.Sprintf("DB error on %s.%s: %v", e.Table, e.Operation, e.Err)
}

func (e *DBError) Unwrap() error { return e.Err }

// ValidationError represents input validation errors
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on %s: %s", e.Field, e.Message)
}
