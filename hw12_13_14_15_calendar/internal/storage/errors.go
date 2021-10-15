package storage

import "fmt"

var (
	ErrStartAtBusy = fmt.Errorf("the start_at time is busy")
	ErrNotFound    = fmt.Errorf("not found")
)
