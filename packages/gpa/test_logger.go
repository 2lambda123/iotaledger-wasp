package gpa

import "github.com/iotaledger/hive.go/ierrors"

// Useful in tests, to make some warnings apparent.
type panicLogger struct{}

func NewPanicLogger() Logger {
	return &panicLogger{}
}

func (*panicLogger) Warnf(msg string, args ...any) {
	panic(ierrors.Errorf(msg, args...))
}
