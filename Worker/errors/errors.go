// errors/errors.go
package errors

import "fmt"

type WorkerError struct {
	Message string
}

func (e *WorkerError) Error() string {
	return fmt.Sprintf("Worker error: %s", e.Message)
}
