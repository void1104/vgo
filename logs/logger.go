package logs

import (
	"io"
	"os"
	"sync"
)

type Logger struct {
	// 一般将日志输出到一个文件，也可以输出到Kafka
	Out io.Writer
	// Flag for whether to log caller info (off by default)
	ReportCaller bool

	// The logging level the logger should log at.
	Level Level
	// Used to sync writing to the log. Locking is enabled by Default
	mu MutexWrap
	// Reusable empty entry
	entryPool sync.Pool
	// Function to exit the application, default to `os.Exit()`
	ExitFunc exitFunc
}

type MutexWrap struct {
	lock     sync.Mutex
	disabled bool
}

func (mw *MutexWrap) Lock() {
	if !mw.disabled {
		mw.lock.Lock()
	}
}

func (mw *MutexWrap) Unlock() {
	if !mw.disabled {
		mw.lock.Unlock()
	}
}

func (mw *MutexWrap) Disable() {
	mw.disabled = true
}

type exitFunc func(int)

// New 建议创建一个全局实例log，也可以自定义通过 &logs.Logger{}自定义生成日志对象
func New() *Logger {
	return &Logger{
		Out:          os.Stderr,
		Level:        InfoLevel,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}
}

func (logger *Logger) newEntry() *Entry {
	// 这里使用池来缓存对象，避免项目大量重复地创建许多对象。
	entry, ok := logger.entryPool.Get().(*Entry)
	if ok {
		return entry
	}
	return NewEntry(logger)
}

func (logger *Logger) releaseEntry(entry *Entry) {
	logger.entryPool.Put(entry)
}

func (logger *Logger) Log(level Level, args ...interface{}) {
	entry := logger.newEntry()
	entry.Log(level, args...)
	logger.releaseEntry(entry)
}

func (logger *Logger) Trace(args ...interface{}) {
	logger.Log(TraceLevel, args...)
}

func (logger *Logger) Debug(args ...interface{}) {
	logger.Log(DebugLevel, args...)
}

func (logger *Logger) Info(args ...interface{}) {
	logger.Log(InfoLevel, args...)
}

func (logger *Logger) Warn(args ...interface{}) {
	logger.Log(WarnLevel, args...)
}

func (logger *Logger) Error(args ...interface{}) {
	logger.Log(ErrorLevel, args...)
}

func (logger *Logger) Panic(args ...interface{}) {
	logger.Log(PanicLevel, args...)
}
