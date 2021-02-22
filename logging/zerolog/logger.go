package zerolog

import (
	zerologpkg "github.com/rs/zerolog"
)

//
// Logger
//

type Logger struct {
	Logger zerologpkg.Logger
}

// logging.Logger interface

func (self *Logger) Critical(message string) {
	self.Logger.Error().Msg(message)
}

func (self *Logger) Criticalf(format string, values ...interface{}) {
	self.Logger.Error().Msgf(format, values...)
}

func (self *Logger) Error(message string) {
	self.Logger.Error().Msg(message)
}

func (self *Logger) Errorf(format string, values ...interface{}) {
	self.Logger.Error().Msgf(format, values...)
}

func (self *Logger) Warning(message string) {
	self.Logger.Warn().Msg(message)
}

func (self *Logger) Warningf(format string, values ...interface{}) {
	self.Logger.Warn().Msgf(format, values...)
}

func (self *Logger) Notice(message string) {
	self.Logger.Info().Msg(message)
}

func (self *Logger) Noticef(format string, values ...interface{}) {
	self.Logger.Info().Msgf(format, values...)
}

func (self *Logger) Info(message string) {
	self.Logger.Debug().Msg(message)
}

func (self *Logger) Infof(format string, values ...interface{}) {
	self.Logger.Debug().Msgf(format, values...)
}

func (self *Logger) Debug(message string) {
	self.Logger.Trace().Msg(message)
}

func (self *Logger) Debugf(format string, values ...interface{}) {
	self.Logger.Trace().Msgf(format, values...)
}
