package response

import "fmt"

type Error struct {
	Code    uint
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}