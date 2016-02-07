package logging

import (
	"io"
	"log"
	"os"
)

// LogLevel the levels supported
type LogLevel int

const (
	// ALL default level - everything gets through
	ALL LogLevel = 0
	// TRACE trace = 10
	TRACE LogLevel = 10
	// DEBUG debug = 20
	DEBUG LogLevel = 20
	// INFO info = 30
	INFO LogLevel = 30
	// WARN warn = 40
	WARN LogLevel = 40
	// ERROR fatal = 50 --> causes a log.Fatal
	ERROR LogLevel = 50
)

// Logger makes it easy to configure the log level and destinations
type Logger struct {
	currentLevel LogLevel
	name         string
	traceHandler *log.Logger
	debugHandler *log.Logger
	infoHandler  *log.Logger
	warnHandler  *log.Logger
	errorHandler *log.Logger
}

var flags = log.Ldate | log.Ltime | log.Lshortfile

// NewLogger builds a new logger with default settings to stderr & stdout
func NewLogger(name string) *Logger {
	logger := new(Logger)

	logger.currentLevel = DEBUG
	logger.name = name
	logger.traceHandler = log.New(os.Stdout, "", flags)
	logger.debugHandler = log.New(os.Stdout, "", flags)
	logger.infoHandler = log.New(os.Stdout, "", flags)
	logger.warnHandler = log.New(os.Stdout, "", flags)
	logger.errorHandler = log.New(os.Stderr, "", flags)

	return logger
}

// SetLevel sets the level of the whole thing
func (logger *Logger) SetLevel(level LogLevel) *Logger {
	logger.currentLevel = level
	return logger
}

// SetWriter overrides the current writer for that level
func (logger *Logger) SetWriter(level LogLevel, writer io.Writer) *Logger {
	switch level {
	case TRACE:
		logger.traceHandler = log.New(writer, "", flags)
	case DEBUG:
		logger.debugHandler = log.New(writer, "", flags)
	case INFO:
		logger.infoHandler = log.New(writer, "", flags)
	case WARN:
		logger.warnHandler = log.New(writer, "", flags)
	case ERROR:
		logger.errorHandler = log.New(writer, "", flags)
	}
	return logger
}

// Trace -- print at level
func (logger Logger) Trace(format string, args ...interface{}) {
	withLevel(logger, TRACE, format, args)
}

// Debug -- print at level
func (logger Logger) Debug(format string, args ...interface{}) {
	withLevel(logger, DEBUG, format, args)
}

// Info -- print at level
func (logger Logger) Info(format string, args ...interface{}) {
	withLevel(logger, INFO, format, args)
}

// Warn -- print at level
func (logger Logger) Warn(format string, args ...interface{}) {
	withLevel(logger, WARN, format, args)
}

// Err -- print at level
func (logger Logger) Err(format string, args ...interface{}) {
	withLevel(logger, ERROR, format, args)
}

func withLevel(logger Logger, level LogLevel, format string, args []interface{}) {
	if logger.currentLevel < level {
		var errStr string
		var l *log.Logger
		switch level {
		case TRACE:
			errStr = "TRACE"
			l = logger.traceHandler
		case DEBUG:
			errStr = "DEBUG"
			l = logger.debugHandler
		case INFO:
			errStr = "INFO"
			l = logger.infoHandler
		case WARN:
			errStr = "WARN"
			l = logger.warnHandler
		case ERROR:
			errStr = "ERROR"
			l = logger.errorHandler
		}

		fmtToUse := errStr + ":" + logger.name + ":" + format + "\n"
		l.Printf(fmtToUse, args...)
	}
}
