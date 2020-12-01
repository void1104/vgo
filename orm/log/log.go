package log

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
)

/**
自定日志库，因为log标准库没有日志分级，不打印文件和行号，
这就意味着我们很难快速知道是哪个地方发生了错误

这个简易的 log 库具备以下特性：
	1. 支持日志分级（Info、Error、Disabled 三级）。
	2. 不同层级日志显示时使用不同的颜色区分。
	3. 显示打印日志代码对应的文件名和行号。
*/

/**
第一步：创建2个日志实例，分别用于打印Info和Error日志
*/
var (
	errorLog = log.New(os.Stdout, "\\033[31m[error]\\033[0m", log.LstdFlags|log.Lshortfile)
	infoLog  = log.New(os.Stdout, "\\033[34m[info ]\\033[0m", log.LstdFlags|log.Lshortfile)
	loggers  = []*log.Logger{errorLog, infoLog}
	mu       sync.Mutex
)

// log methods
var (
	Error = errorLog.Println
	Errof = errorLog.Printf
	Info  = infoLog.Println
	infof = infoLog.Printf
)

/**
第二步：支持设置日志的层级（InfoLevel, ErrorLevel, Disabled）
*/
// log level
const (
	InfoLevel = iota
	ErrorLevel
	Disabled
)

// SetLevel controls log level
func SetLevel(level int) {
	mu.Lock()
	defer mu.Unlock()

	for _, logger := range loggers {
		logger.SetOutput(os.Stdout)
	}

	if ErrorLevel < level {
		errorLog.SetOutput(ioutil.Discard)
	}

	if InfoLevel < level {
		infoLog.SetOutput(ioutil.Discard)
	}
}
