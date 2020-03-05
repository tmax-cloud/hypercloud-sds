package tests

import "time"

const (
	PollingIntervalDefault              = 3 * time.Second
	TimeoutForDeletingNamespace         = 500 * time.Second
	PollingIntervalForDeletingNamespace = 10 * time.Second
)
