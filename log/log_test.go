package log

import "testing"

func TestTrace(t *testing.T) {
	Trace("trace testing")
}

func TestDebug(t *testing.T) {
	Debug("debug testing")
}

func TestInfo(t *testing.T) {
	Info("info testing")
}

func TestWarn(t *testing.T) {
	Warn("warn testing")
}

func TestError(t *testing.T) {
	Error("error testing")
}

func TestPanic(t *testing.T) {
	Panic("panic testing")
}