package auth

import "fmt"

func NewError(tag string, message string, params ...interface{}) FieldError {
	return FieldError{
		Tag:     tag,
		Message: message,
		Params:  params,
	}
}

func (fe FieldError) Error() string {
	return fmt.Sprintf("message: %s, tag: %s, params: %v", fe.Message, fe.Tag, fe.Params)
}
