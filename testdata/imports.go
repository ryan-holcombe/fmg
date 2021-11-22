package testdata

import (
	"net/url"
	"time"
)

// Imports struct to help test imports
type Imports struct {
	ID          int64 `fmgen:"-"`
	Name        string
	Age         int64 `fmgen:"optional"`
	LastUpdated time.Time
	Duration    time.Duration
	BaseURL     url.URL
}
