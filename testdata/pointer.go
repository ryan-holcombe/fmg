package testdata

import "time"

// Pointer struct to help test pointers
type Pointer struct {
	ID          int64 `fmgen:"-"`
	Name        string
	Age         int64 `fmgen:"optional"`
	PtrS        *string
	PtrOpt      *string `fmgen:"optional"`
	PtrI        *int
	LastUpdated *time.Time
}
