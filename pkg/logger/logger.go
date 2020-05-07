// Package logger wraps the zap logger from uber and adds fluentd-like
// formatting. This gives us a fast, structured logger which is also supported
// by Google Cloud Platform and StackDriver.
package logger

import (
	"runtime/debug"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New() *zap.Logger {
	config := config(true)

	logger, err := config.Build(wrap())
	if err != nil {
		panic(err)
	}

	return logger
}

// config will adjust the zap defaults to behave better on StackDriver
func config(debug bool) zap.Config {
	config := zap.NewProductionConfig()

	if debug {
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		config.Development = true
	}

	config.EncoderConfig.LevelKey = "severity"
	config.EncoderConfig.TimeKey = "time"
	config.EncoderConfig.MessageKey = "message"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stdout"}
	config.DisableStacktrace = true

	return config
}

type core struct {
	zapcore.Core
}

// Check overrides the embedded zapcore.Core Check method on our `core` struct
func (c *core) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}

	return ce
}

// add stacktraces to the log message on any statement with a
// level of error or above. This is a workaround to make tracking errors
// better on StackDriver because zap doesn't allow us to split up the
// location of the error into multiple values:
// https://github.com/uber-go/zap/issues/627
func (c *core) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	if ent.Level >= zapcore.ErrorLevel {
		for _, f := range fields {
			if f.Type == zapcore.ErrorType {
				err, ok := f.Interface.(error)
				if ok {
					ent.Message += ": " + err.Error()
				}
			}
		}

		ent.Message += "\n" + string(debug.Stack())
	}

	return c.Core.Write(ent, fields)
}

func (c *core) With(fields []zap.Field) zapcore.Core {
	return &core{c.Core.With(fields)}
}

func wrap() zap.Option {
	return zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return &core{c}
	})
}
