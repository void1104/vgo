package log

import (
	"fmt"
	"os"
	"time"
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
	//Context core.Context

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

	// 将日志输入到指定文件路径
	file, err := os.OpenFile(entry.Logger.LogPath, os.O_WRONLY, 0644)
	if err != nil && os.IsNotExist(err) {
		file, err = os.Create(entry.Logger.LogPath)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "fail open local log file")
			os.Exit(1)
		}
	}
	defer file.Close()

	// file 类型实现了io.Writer
	n, _ := file.Seek(0, 2)
	_, fErr := file.WriteAt(log, n)
	if fErr != nil {
		_, _ = fmt.Fprintf(os.Stderr, "fail write log to local file")
		os.Exit(1)
	}

	if _, err = entry.Logger.Out.Write(log); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to write to log, #{err}\n")
	}
}

func (entry *Entry) Log(level Level, args ...interface{}) {
	entry.log(level, fmt.Sprint(args...))
}

func (entry *Entry) formatter() []byte {
	timestamp := "[" + time.Now().Format("2006-01-02 15:04:05") + "]"
	level := "[" + entry.Level.String() + "] "
	msg := entry.Message
	return []byte(timestamp + level + msg + "\n")
}
