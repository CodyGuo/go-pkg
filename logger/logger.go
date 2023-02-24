package logger

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/CodyGuo/go-pkg/fs"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	TimeFormat = "2006-01-02T15:04:05.000"
)

var (
	_logger       Logger
	_accessLogger Logger
)

type Logger struct {
	once sync.Once
	skip int

	fileLogger    zerolog.Logger
	consoleLogger zerolog.Logger
	rollingLogger *lumberjack.Logger
}

func New(c Config) (logger Logger, err error) {
	if c.EnableFile {
		if err = fs.MkdirAll(fs.Dir(c.FilePath), os.ModePerm); err != nil {
			return
		}
	}

	level := levelToZeroLevel(c.Level)
	timeFormat := c.TimeFormat
	if c.TimeFormat == "" {
		timeFormat = TimeFormat
	}
	SetLogTime(timeFormat, c.UTCTime)

	consoleLogger := zerolog.Nop()
	if c.EnableConsole {
		consoleOutput := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: timeFormat,
		}
		consoleLogger = zerolog.New(consoleOutput).Level(level)
	}

	fileLogger := zerolog.Nop()
	rollingLogger := &lumberjack.Logger{}
	if c.EnableFile {
		rollingLogger = &lumberjack.Logger{
			Filename:   c.FilePath,
			MaxSize:    c.MaxSize,
			MaxAge:     c.MaxAge,
			MaxBackups: c.MaxBackups,
			LocalTime:  !c.UTCTime,
			Compress:   c.Compress,
		}
		fileLogger = zerolog.New(rollingLogger).Level(level)
	}

	logger = Logger{
		once:          sync.Once{},
		skip:          4,
		fileLogger:    fileLogger,
		consoleLogger: consoleLogger,
		rollingLogger: rollingLogger,
	}
	return
}

func (l Logger) Debug(format string, v ...any) {
	l.log(DEBUG, format, v...)
}

func (l Logger) Info(format string, v ...any) {
	l.log(INFO, format, v...)
}

func (l Logger) Warn(format string, v ...any) {
	l.log(WARN, format, v...)
}

func (l Logger) Error(format string, v ...any) {
	l.log(ERROR, format, v...)
}

func (l Logger) DebugToFile(format string, v ...any) {
	l.logToFile(DEBUG, format, v...)
}

func (l Logger) InfoToFile(format string, v ...any) {
	l.logToFile(INFO, format, v...)
}

func (l Logger) WarnToFile(format string, v ...any) {
	l.logToFile(WARN, format, v...)
}

func (l Logger) ErrorToFile(format string, v ...any) {
	l.logToFile(ERROR, format, v...)
}

func (l Logger) DebugToConsole(format string, v ...any) {
	l.logToConsole(DEBUG, format, v...)
}

func (l Logger) InfoToConsole(format string, v ...any) {
	l.logToConsole(INFO, format, v...)
}

func (l Logger) WarnToConsole(format string, v ...any) {
	l.logToConsole(WARN, format, v...)
}

func (l Logger) ErrorToConsole(format string, v ...any) {
	l.logToConsole(ERROR, format, v...)
}

func (l Logger) With(key string, val any) *Logger {
	return l.updateLoggerContext(key, val)
}

func (l Logger) WithSender(sender string) *Logger {
	return l.updateLoggerContext("sender", sender)
}

func (l Logger) WithRequestID(id string) *Logger {
	return l.updateLoggerContext("request_id", id)
}

func (l Logger) AutoSkipFrameCount() *Logger {
	l.once.Do(func() {
		l.skip += 1
	})
	return &l
}

func (l Logger) WithSkipFrameCount(skip int) *Logger {
	l.skip += skip
	return &l
}

func (l Logger) WithCaller() *Logger {
	l.fileLogger = l.fileLogger.With().CallerWithSkipFrameCount(l.skip).Logger()
	l.consoleLogger = l.consoleLogger.With().CallerWithSkipFrameCount(l.skip).Logger()
	return &l
}

func (l Logger) WithHook(h Hook) *Logger {
	l.fileLogger = l.fileLogger.Hook(ZeroHook(h))
	l.consoleLogger = l.consoleLogger.Hook(ZeroHook(h))
	return &l
}

func (l Logger) WithHookFunc(h HookFunc) *Logger {
	l.fileLogger = l.fileLogger.Hook(ZeroHook(h))
	l.consoleLogger = l.consoleLogger.Hook(ZeroHook(h))
	return &l
}

func (l Logger) updateLoggerContext(key string, val any) *Logger {
	var loggers = []*zerolog.Logger{&l.fileLogger, &l.consoleLogger}

	for _, logger := range loggers {
		var context zerolog.Context
		switch v := val.(type) {
		case string:
			context = logger.With().Str(key, v)
		case []string:
			context = logger.With().Strs(key, v)
		case fmt.Stringer:
			context = logger.With().Stringer(key, v)
		case []byte:
			context = logger.With().Bytes(key, v)
		case error:
			context = logger.With().AnErr(key, v)
		case []error:
			context = logger.With().Errs(key, v)
		case bool:
			context = logger.With().Bool(key, v)
		case []bool:
			context = logger.With().Bools(key, v)
		case int:
			context = logger.With().Int(key, v)
		case []int:
			context = logger.With().Ints(key, v)
		case int8:
			context = logger.With().Int8(key, v)
		case []int8:
			context = logger.With().Ints8(key, v)
		case int16:
			context = logger.With().Int16(key, v)
		case []int16:
			context = logger.With().Ints16(key, v)
		case int32:
			context = logger.With().Int32(key, v)
		case []int32:
			context = logger.With().Ints32(key, v)
		case int64:
			context = logger.With().Int64(key, v)
		case []int64:
			context = logger.With().Ints64(key, v)
		case uint:
			context = logger.With().Uint(key, v)
		case []uint:
			context = logger.With().Uints(key, v)
		case uint8:
			context = logger.With().Uint8(key, v)
		// case []uint8: // []byte
		// 	context = logger.With().Uints8(key, v)
		case uint16:
			context = logger.With().Uint16(key, v)
		case []uint16:
			context = logger.With().Uints16(key, v)
		case uint32:
			context = logger.With().Uint32(key, v)
		case []uint32:
			context = logger.With().Uints32(key, v)
		case uint64:
			context = logger.With().Uint64(key, v)
		case []uint64:
			context = logger.With().Uints64(key, v)
		case float32:
			context = logger.With().Float32(key, v)
		case []float32:
			context = logger.With().Floats32(key, v)
		case float64:
			context = logger.With().Float64(key, v)
		case []float64:
			context = logger.With().Floats64(key, v)
		case time.Time:
			context = logger.With().Time(key, v)
		case []time.Time:
			context = logger.With().Times(key, v)
		case time.Duration:
			context = logger.With().Dur(key, v)
		case []time.Duration:
			context = logger.With().Durs(key, v)
		case net.IP:
			context = logger.With().IPAddr(key, v)
		case net.IPNet:
			context = logger.With().IPPrefix(key, v)
		case net.HardwareAddr:
			context = logger.With().MACAddr(key, v)
		default:
			context = logger.With().Interface(key, v)
		}
		*logger = context.Logger()
	}

	return &l
}

func (l Logger) log(level Level, format string, v ...any) {
	l.levelEvent(level, l.fileLogger).Timestamp().Msgf(format, v...)
	l.levelEvent(level, l.consoleLogger).Timestamp().Msgf(format, v...)
}

func (l Logger) logToFile(level Level, format string, v ...any) {
	l.levelEvent(level, l.fileLogger).Timestamp().Msgf(format, v...)
}

func (l Logger) logToConsole(level Level, format string, v ...any) {
	l.levelEvent(level, l.consoleLogger).Timestamp().Msgf(format, v...)
}

func (l Logger) levelEvent(level Level, logger zerolog.Logger) *zerolog.Event {
	var ev *zerolog.Event
	switch level {
	case DEBUG:
		ev = logger.Debug()
	case INFO:
		ev = logger.Info()
	case WARN:
		ev = logger.Warn()
	case ERROR:
		ev = logger.Error()
	default:
		ev = logger.Info()
	}
	return ev
}

func GetLogger() *Logger {
	return _logger.AutoSkipFrameCount()
}

func GetAccessLogger() *Logger {
	return &_accessLogger
}

func With(key string, a any) *Logger {
	return WithSkipFrameCount(-1).With(key, a)
}

func WithSender(sender string) *Logger {
	return WithSkipFrameCount(-1).WithSender(sender)
}

func WithRequestID(id string) *Logger {
	return WithSkipFrameCount(-1).WithRequestID(id)
}

func AutoSkipFrameCount() *Logger {
	return WithSkipFrameCount(-1).AutoSkipFrameCount()
}

func WithSkipFrameCount(skip int) *Logger {
	return GetLogger().WithSkipFrameCount(skip)
}

func WithCaller() *Logger {
	return WithSkipFrameCount(-1).WithCaller()
}

func WithHook(h Hook) *Logger {
	return WithSkipFrameCount(-1).WithHook(h)
}

func WithHookFunc(h HookFunc) *Logger {
	return WithSkipFrameCount(-1).WithHookFunc(h)
}

func Debug(format string, v ...any) {
	GetLogger().Debug(format, v...)
}

func Info(format string, v ...any) {
	GetLogger().Info(format, v...)
}

func Warn(format string, v ...any) {
	GetLogger().Warn(format, v...)
}

func Error(format string, v ...any) {
	GetLogger().Error(format, v...)
}

func DebugToFile(format string, v ...any) {
	GetLogger().DebugToFile(format, v...)
}

func InfoToFile(format string, v ...any) {
	GetLogger().InfoToFile(format, v...)
}

func WarnToFile(format string, v ...any) {
	GetLogger().WarnToFile(format, v...)
}

func ErrorToFile(format string, v ...any) {
	GetLogger().ErrorToFile(format, v...)
}

func DebugToConsole(format string, v ...any) {
	GetLogger().DebugToConsole(format, v...)
}

func InfoToConsole(format string, v ...any) {
	GetLogger().InfoToConsole(format, v...)
}

func WarnToConsole(format string, v ...any) {
	GetLogger().WarnToConsole(format, v...)
}

func ErrorToConsole(format string, v ...any) {
	GetLogger().ErrorToConsole(format, v...)
}
