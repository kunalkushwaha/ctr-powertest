package testcase

import (
	"context"
)

const (
	testImage         = "docker.io/library/alpine:latest"
	testContainerName = "test-powertest"
)

// Testcases interface to implement testcases
type Testcases interface {
	RunAllTests(context.Context) error
}
