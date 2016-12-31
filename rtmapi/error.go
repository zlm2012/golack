package rtmapi

// MalformedPayloadError represents an error that given JSON payload is not properly formatted.
// e.g. required fields are not given, or payload is not a valid JSON string.
type MalformedPayloadError struct {
	Err string
}

// NewMalformedPayloadError creates new MalformedPayloadError instance with given arguments.
func NewMalformedPayloadError(str string) *MalformedPayloadError {
	return &MalformedPayloadError{Err: str}
}

// Error returns its error string.
func (e *MalformedPayloadError) Error() string {
	return e.Err
}
