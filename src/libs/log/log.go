package log

import (
	"github.com/astaxie/beego/logs"
	"libs/lua"
	"strings"
)

var logger *logs.BeeLogger

func init() {

	L, err := lua.NewLua("conf/app.lua")
	if err != nil {
		panic(err)
	}

	levelMap := make(map[string]int)
	levelMap["Trace"] = logs.LevelTrace
	levelMap["Debug"] = logs.LevelDebug
	levelMap["Info"] = logs.LevelInfo
	levelMap["Warn"] = logs.LevelWarn
	levelMap["Error"] = logs.LevelError
	levelMap["Critical"] = logs.LevelCritical

	debugLevel, _ := levelMap[strings.Title(L.GetString("level"))]

	L.Close()

	// init focust core
	logger = logs.NewLogger(1000)
	logger.SetLogger("console", "")
	logger.SetLevel(debugLevel)
}

// SetLogLevel sets the global log level used by the simple logger.
func SetLevel(l int) {
	logger.SetLevel(l)
}

// Trace logs a message at trace level.
func Trace(format string, v ...interface{}) {
	logger.Trace(format, v...)
}

// Debug logs a message at debug level.
func Debug(format string, v ...interface{}) {
	logger.Debug(format, v...)
}

// Info logs a message at info level.
func Info(format string, v ...interface{}) {
	logger.Info(format, v...)
}

// Warning logs a message at warning level.
func Warn(format string, v ...interface{}) {
	logger.Warn(format, v...)
}

// Error logs a message at error level.
func Error(format string, v ...interface{}) {
	logger.Error(format, v...)
}

// Critical logs a message at critical level.
func Critical(format string, v ...interface{}) {
	logger.Critical(format, v...)
}
