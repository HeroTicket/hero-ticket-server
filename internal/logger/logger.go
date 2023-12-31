package logger

import (
	"go.uber.org/zap"
)

var (
	l *zap.Logger
	s *zap.SugaredLogger
)

func New(development bool, withArgs ...interface{}) error {
	var err error

	if development {
		l, err = zap.NewDevelopment(zap.AddCallerSkip(1))
	} else {
		l, err = zap.NewProduction(zap.AddCallerSkip(1))
	}
	if err != nil {
		return err
	}

	s = l.Sugar()

	s = s.With(withArgs...)

	return nil
}

func Sync() error {
	if l == nil {
		return nil
	}

	return l.Sync()
}

func Info(msg string, args ...interface{}) {
	if s == nil {
		zap.L().Sugar().Infow(msg, args...)
		return
	}
	s.Infow(msg, args...)
}

func Debug(msg string, args ...interface{}) {
	if s == nil {
		zap.L().Sugar().Debugw(msg, args...)
		return
	}
	s.Debugw(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	if s == nil {
		zap.L().Sugar().Warnw(msg, args...)
		return
	}
	s.Warnw(msg, args...)
}

func Error(msg string, args ...interface{}) {
	if s == nil {
		zap.L().Sugar().Errorw(msg, args...)
		return
	}
	s.Errorw(msg, args...)
}

func Panic(msg string, args ...interface{}) {
	if s == nil {
		zap.L().Sugar().Panicw(msg, args...)
		return
	}
	s.Panicw(msg, args...)
}
