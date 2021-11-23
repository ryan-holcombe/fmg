package testdata

import "time"

// Sample simple struct fmgen:omit
type Sample struct {
	ID          int64 `fmgen:"-"`
	Name        string
	Age         int64 `fmgen:"optional"`
	LastUpdated time.Time
}

// Simple struct with just a name
type Simple struct {
	Name string
}
