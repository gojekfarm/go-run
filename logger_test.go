package go_run

import "fmt"

type TestLogger struct{}

func (l TestLogger) Debugf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (l TestLogger) Infof(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (l TestLogger) Warnf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (l TestLogger) Errorf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (l TestLogger) Fatalf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (l TestLogger) Panicf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}
