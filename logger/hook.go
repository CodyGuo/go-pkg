package logger

import "github.com/rs/zerolog"

type Event = zerolog.Event

type Hook interface {
	Run(e *Event, level Level, message string)
}

// HookFunc is an adaptor to allow the use of an ordinary function
// as a Hook.
type HookFunc func(e *Event, level Level, message string)

// Run implements the zerolog.Hook interface.
func (h HookFunc) Run(e *Event, level Level, message string) {
	h(e, level, message)
}

type zeroHook struct {
	Hook
}

func (h zeroHook) Run(e *zerolog.Event, level zerolog.Level, message string) {
	h.Hook.Run(e, zeroLevelToLevel(level), message)
}

func ZeroHook(h Hook) zeroHook {
	return zeroHook{h}
}
