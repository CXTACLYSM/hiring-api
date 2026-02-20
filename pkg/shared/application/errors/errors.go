package errors

type ApplicationError struct {
	Message string
	Err     error
}

func (e *ApplicationError) Error() string {
	return e.Message
}
