package errors

const MessageInternalServerError = "internal server error"

type AppError struct {
	Message    string
	StatusCode int
	Err        error
}

func (e *AppError) Error() string {
	return e.Message
}
