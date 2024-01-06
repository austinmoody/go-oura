package go_oura

import "fmt"

type OuraError struct {
	Code    int
	Message string
}

func (e *OuraError) Error() string {
	return fmt.Sprintf("OuraError %d - %s", e.Code, e.Message)
}
