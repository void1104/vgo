package logs

import (
	"fmt"
	"os"
	"time"
	"vgo/context"
)

//var (
//	// qualified package name, cached at first use
//	logPackage string
//
//	// Positions in the call stack when tracing to report the calling method
//	minimumCallerDepth int
//
//	// Used for caller information initialisation
//	callerInitOnce sync.Once
//)

//const (
//	maximumCallerDepth int = 25
//	knownLogFrames     int = 4
//)

//func init() {
//	// start at the bottom of the stack befor the package-name cache is primed
//	minimumCallerDepth = 1
//}

type Entry struct {
	Logger *Logger

	// Time at which the log entry was created
	Time time.Time

	// Level the log entry level
	Level Level

	// Caller method, with package name
	//Caller *runtime.Frame

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
	// TODO 日志内容修改
	log := time.Now().String() + "hello 日志测试"

	entry.Logger.mu.Lock()
	defer entry.Logger.mu.Unlock()
	// TODO 日志文件输入路径
	if _, err := entry.Logger.Out.Write([]byte(log)); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to log, #{err}\n")
	}
}

func (entry *Entry) Log(level Level, args ...interface{}) {
	entry.log(level, fmt.Sprint(args...))
}
