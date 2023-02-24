package logger

import "github.com/rs/zerolog"

const (
	DEBUG Level = iota - 1
	INFO
	WARN
	ERROR
)

type Level int

func (l Level) String() string {
	var level string
	switch l {
	case DEBUG:
		level = "debug"
	case INFO:
		level = "info"
	case WARN:
		level = "warn"
	case ERROR:
		level = "error"
	default:
		level = "unknown"
	}
	return level
}

func levelToZeroLevel(level string) zerolog.Level {
	switch level {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	default:
		return zerolog.InfoLevel
	}
}

func zeroLevelToLevel(level zerolog.Level) Level {
	switch level {
	case zerolog.DebugLevel:
		return DEBUG
	case zerolog.InfoLevel:
		return INFO
	case zerolog.WarnLevel:
		return WARN
	case zerolog.ErrorLevel:
		return ERROR
	default:
		return INFO
	}
}
