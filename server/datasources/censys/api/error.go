package api

import "fmt"

type Error struct {
	ErrorCode   int    `json:"error_code"`
	ErrorString string `json:"error"`
}

func (de *Error) Error() string {
	return fmt.Sprintf("%s (%d)", de.ErrorString, de.ErrorCode)
}
