package errors

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func New(code int, message string) *Error {
	return &Error{Code: code, Message: message}
}

func (e *Error) Error() string {
	return e.Message
}
