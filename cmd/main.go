package main

import (
	"github.com/sebastianreh/compass-code-assessment/internal"
	"time"
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

func main() {
	// Build dependencies
	build := internal.Build()

	// Generate start timestamp
	start := time.Now()
	build.Logger.Infof("Start processing at %v", start.Format(timeFormat))
	_, err := build.Service.Evaluate()
	if err != nil {
		panic(err)
	}

	// Log finish time
	build.Logger.Infof("Finish processing at %v", time.Now().Format(timeFormat))
	build.Logger.Infof("Process took %v", time.Since(start))
}
