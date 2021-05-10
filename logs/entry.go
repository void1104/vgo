package logs

import (
	"fmt"
	"os"
	"time"
	"vgo/context"
)

type Entry struct {
	Logger *Logger

	// Time at which the log entry was created
	Time time.Time

	// Level the log entry level
	Level Level

	// Message passed to Trace, Debug, Info, Warn, Error, Fatal or Panic
	Message string

	// Contains the context set by the user. Useful for hook processing etc.
	Context context.Context

	// err may contain a field formatting error
	err string
}

func NewEntry(logger *Logger) *Entry {
	return &Entry{
		Logger: logger,
	}
}

func (entry *Entry) log(level Level, msg string) {

	if entry.Time.IsZero() {
		entry.Time = time.Now()
	}
	entry.Level = level
	entry.Message = msg

	entry.write()
}

func (entry *Entry) write() {

	log := entry.formatter()

	entry.Logger.mu.Lock()
	defer entry.Logger.mu.Unlock()
	// TODO 日志文件输入路径
	if _, err := entry.Logger.Out.Write(log); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to log, #{err}\n")
	}
}

func (entry *Entry) Log(level Level, args ...interface{}) {
	entry.log(level, fmt.Sprint(args...))
}

func (entry *Entry) formatter() []byte {
	timestamp := "[" + entry.Time.Format("2006-01-02 15:03:04") + "]"
	level := "[" + entry.Level.String() + "] "
	msg := entry.Message
	return []byte(timestamp + level + msg)
}
