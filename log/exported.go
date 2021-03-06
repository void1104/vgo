package log

import "io"

var (
	// std is the name of the standard logger in stdlib `log`
	std = newLogger()
)

// Trace logs a message at level Trace on the standard logger.
func Trace(args ...interface{}) {
	std.Trace(args...)
}

// Debug logs a message at level Debug on the standard logger
func Debug(args ...interface{}) {
	std.Debug(args...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	std.Info(args...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	std.Warn(args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	std.Error(args...)
}

// Panic logs a message at level Panic on the standard logger.
func Panic(args ...interface{}) {
	std.Panic(args...)
}

// SetLogPath set the file path where the log to writer.
func SetLogPath(logPath string) {
	std.SetLogPath(logPath)
}

// SetOutput set the file output
func SetOutput(output io.Writer) {
	std.SetOutput(output)
}

////
//func AddHook() {
//
//}
