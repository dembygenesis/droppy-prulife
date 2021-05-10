package utils

// ApplicationError returns more detailed information about the error
type ApplicationError struct {
	HttpStatus int
	Message    string
	Error      interface{}
}
