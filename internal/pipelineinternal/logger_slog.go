package pipelineinternal

import "log/slog"

type slogLogger struct{ l *slog.Logger }

func FromSlog(l *slog.Logger) Logger {
	if l == nil {
		return nil
	}
	return slogLogger{l: l}
}

func (s slogLogger) Debug(msg string, args ...any) { s.l.Debug(msg, args...) }
func (s slogLogger) Info(msg string, args ...any)  { s.l.Info(msg, args...) }
func (s slogLogger) Warn(msg string, args ...any)  { s.l.Warn(msg, args...) }
func (s slogLogger) Error(msg string, args ...any) { s.l.Error(msg, args...) }
