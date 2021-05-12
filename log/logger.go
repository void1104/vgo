package log

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

// 各日志级别
const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

const DefaultLogPath = "./log.txt"

// Level type
type Level uint32

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
	// log file path
	LogPath string
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

// NewLogger 建议创建一个全局实例log，也可以自定义通过 &logs.Logger{}自定义生成日志对象
func newLogger() *Logger {
	return &Logger{
		Out:          os.Stdout,
		Level:        InfoLevel,
		ExitFunc:     os.Exit,
		ReportCaller: false,
		LogPath:      DefaultLogPath,
	}
}

// newEntry 这里使用池来缓存对象，避免项目大量重复地创建许多对象。
func (logger *Logger) newEntry() *Entry {
	entry, ok := logger.entryPool.Get().(*Entry)
	if ok {
		return entry
	}
	return NewEntry(logger)
}

// releaseEntry 将entry对象放回对象池
func (logger *Logger) releaseEntry(entry *Entry) {
	logger.entryPool.Put(entry)
}

// SetLogPath 设置日志输出文件路径
func (logger *Logger) SetLogPath(logPath string) {
	logger.LogPath = logPath
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

// Convert the Level to a string. E.g. PanicLevel becomes "panic".
func (level Level) String() string {
	switch level {
	case TraceLevel:
		return "trace"
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warning"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	case PanicLevel:
		return "panic"
	}
	return "unknown"
}

// ParseLevel takes a string level and returns the log level constant.
func ParseLevel(lvl string) (Level, error) {
	switch strings.ToLower(lvl) {
	case "panic":
		return PanicLevel, nil
	case "fatal":
		return FatalLevel, nil
	case "error":
		return ErrorLevel, nil
	case "warn", "warning":
		return WarnLevel, nil
	case "info":
		return InfoLevel, nil
	case "debug":
		return DebugLevel, nil
	case "trace":
		return TraceLevel, nil
	}

	var l Level
	return l, fmt.Errorf("not a valid logrus Level: %q", lvl)
}
