/* The MIT License (MIT)
Copyright © 2018 by Atlas Lee(atlas@fpay.io)

Permission is hereby granted, free of charge, to any person obtaining a
copy of this software and associated documentation files (the “Software”),
to deal in the Software without restriction, including without limitation
the rights to use, copy, modify, merge, publish, distribute, sublicense,
and/or sell copies of the Software, and to permit persons to whom the
Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
DEALINGS IN THE SOFTWARE.
*/

package zlog /* 格式化日志工具 */

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"sync"
)

const (
	VERBOSE uint8 = iota /* 最详细的输出 */
	TRACE                /* 调试信息 */
	DEBUG                /* 函数进入退出信息 */
	INFO                 /* 服务启动关闭信息及请求访问信息 */
	WARNING              /* 异常警告信息(无需人工干预) */
	ERROR                /* 错误信息(不影响服务，需要人工干预) */
	FATAL                /* 崩溃信息(无法继续提供服务) */
	SILENCE
)

var LogLevelNames [8]string = [8]string{"VERBOSE", "TRACE", "DEBUG", "INFO", "WARNING", "ERROR", "FATAL", "SILENCE"}

var (
	globalLevel  uint8            = VERBOSE                /* 全局日志级别 */
	loggerLevels map[string]uint8 = make(map[string]uint8) /* 指定标志日志级别 */
	mu           sync.Mutex                                /* 全局锁，保证zlog线程安全 */
)

/* 设置全局日志输出级别，低于该级别的日志不会输出 */
func SetLevel(level uint8) {
	globalLevel = level
}

/* 指定具体标志的日志级别，应小于全局级别 */
/* 结合SetLevel，可以只输出指定标志的日志 */
func SetTagLevel(level uint8, tags ...string) {
	mu.Lock()
	for _, tag := range tags {
		loggerLevels[tag] = level
	}
	mu.Unlock()
}

func logf(level uint8, format string, v ...interface{}) {
	callers := make([]uintptr, 1)
	runtime.Callers(4, callers)
	caller := runtime.FuncForPC(callers[0])
	peices := strings.Split(caller.Name(), ".")
	size := len(peices)
	pkg := strings.Join(peices[:size-1], "/")
	tagLevel, ok := loggerLevels[pkg]
	if !ok {
		tagLevel = globalLevel
	}

	if level >= tagLevel {
		method := peices[size-1]
		switch level {
		case VERBOSE, TRACE, DEBUG:
			log.Printf(fmt.Sprintf("[%c[1;32m%s%c[0m][%s: %s] %s", 0x1B, LogLevelNames[level], 0x1B, peices[size-2], method, fmt.Sprintf(format, v...)))
		case INFO, WARNING:
			log.Printf(fmt.Sprintf("[%c[1;37m%s%c[0m][%s: %s] %s", 0x1B, LogLevelNames[level], 0x1B, peices[size-2], method, fmt.Sprintf(format, v...)))
		case ERROR, FATAL:
			log.Printf(fmt.Sprintf("[%c[1;31m%s%c[0m][%s: %s] %s", 0x1B, LogLevelNames[level], 0x1B, peices[size-2], method, fmt.Sprintf(format, v...)))
		default:
			log.Printf(fmt.Sprintf("[%s][%s: %s] %s", LogLevelNames[level], peices[size-2], method, fmt.Sprintf(format, v...)))
		}
	}
}

func logln(level uint8, v ...interface{}) {
	callers := make([]uintptr, 1)
	runtime.Callers(4, callers)
	caller := runtime.FuncForPC(callers[0])
	peices := strings.Split(caller.Name(), ".")
	size := len(peices)
	pkg := strings.Join(peices[:size-1], "/")

	tagLevel, ok := loggerLevels[pkg]
	if !ok {
		tagLevel = globalLevel
	}

	if level >= tagLevel {
		method := peices[size-1]
		switch level {
		case VERBOSE, TRACE, DEBUG:
			log.Printf(fmt.Sprintf("[%c[1;32m%s%c[0m][%s: %s] %s", 0x1B, LogLevelNames[level], 0x1B, peices[size-2], method, fmt.Sprintln(v...)))
		case INFO, WARNING:
			log.Printf(fmt.Sprintf("[%c[1;37m%s%c[0m][%s: %s] %s", 0x1B, LogLevelNames[level], 0x1B, peices[size-2], method, fmt.Sprintln(v...)))
		case ERROR, FATAL:
			log.Printf(fmt.Sprintf("[%c[1;31m%s%c[0m][%s: %s] %s", 0x1B, LogLevelNames[level], 0x1B, peices[size-2], method, fmt.Sprintln(v...)))
		default:
			log.Printf(fmt.Sprintf("[%s][%s: %s] %s", LogLevelNames[level], peices[size-2], method, fmt.Sprintln(v...)))
		}
	}
}

func Logf(level uint8, format string, v ...interface{}) {
	logf(level, format, v...)
}

func Logln(level uint8, v ...interface{}) {
	logln(level, v...)
}

func Verbosef(format string, v ...interface{}) {
	Logf(VERBOSE, format, v...)
}

func Verboseln(v ...interface{}) {
	Logln(VERBOSE, v...)
}

func Tracef(format string, v ...interface{}) {
	Logf(TRACE, format, v...)
}

func Traceln(v ...interface{}) {
	Logln(TRACE, v...)
}

func Debugf(format string, v ...interface{}) {
	Logf(DEBUG, format, v...)
}

func Debugln(v ...interface{}) {
	Logln(DEBUG, v...)
}

func Infof(format string, v ...interface{}) {
	Logf(INFO, format, v...)
}

func Infoln(v ...interface{}) {
	Logln(INFO, v...)
}

func Warningf(format string, v ...interface{}) {
	Logf(WARNING, format, v...)
}

func Warningln(v ...interface{}) {
	Logln(WARNING, v...)
}

func Errorf(format string, v ...interface{}) {
	Logf(ERROR, format, v...)
}

func Errorln(v ...interface{}) {
	Logln(ERROR, v...)
}

func Fatalf(format string, v ...interface{}) {
	Logf(FATAL, format, v...)
}

func Fatalln(v ...interface{}) {
	Logln(FATAL, v...)
}
