package logging

import (
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

type ColorfulLogger struct {
	*logrus.Logger
	colors map[logrus.Level]*color.Color
}

func NewColorfulLogger() *ColorfulLogger {
	logger := &ColorfulLogger{
		Logger: logrus.New(),
		colors: map[logrus.Level]*color.Color{
			logrus.DebugLevel: color.New(color.FgCyan),
			logrus.InfoLevel:  color.New(color.FgGreen),
			logrus.WarnLevel:  color.New(color.FgYellow),
			logrus.ErrorLevel: color.New(color.FgRed),
			logrus.FatalLevel: color.New(color.FgMagenta),
			logrus.PanicLevel: color.New(color.BgRed, color.FgWhite),
		},
	}

	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)

	return logger
}

func (cl *ColorfulLogger) SetLevel(level logrus.Level) {
	cl.Logger.SetLevel(level)
}

func (cl *ColorfulLogger) SetOutput(output io.Writer) {
	cl.Logger.SetOutput(output)
}

func (cl *ColorfulLogger) Debug(args ...interface{}) {
	cl.colors[logrus.DebugLevel].Println(args...)
	cl.Logger.Debug(args...)
}

func (cl *ColorfulLogger) Info(args ...interface{}) {
	cl.colors[logrus.InfoLevel].Println(args...)
	cl.Logger.Info(args...)
}

func (cl *ColorfulLogger) Warn(args ...interface{}) {
	cl.colors[logrus.WarnLevel].Println(args...)
	cl.Logger.Warn(args...)
}

func (cl *ColorfulLogger) Error(args ...interface{}) {
	cl.colors[logrus.ErrorLevel].Println(args...)
	cl.Logger.Error(args...)
}

func (cl *ColorfulLogger) Fatal(args ...interface{}) {
	cl.colors[logrus.FatalLevel].Println(args...)
	cl.Logger.Fatal(args...)
}

func (cl *ColorfulLogger) Panic(args ...interface{}) {
	cl.colors[logrus.PanicLevel].Println(args...)
	cl.Logger.Panic(args...)
}

func (cl *ColorfulLogger) Debugf(format string, args ...interface{}) {
	cl.colors[logrus.DebugLevel].Printf(format, args...)
	cl.Logger.Debugf(format, args...)
}

func (cl *ColorfulLogger) Infof(format string, args ...interface{}) {
	cl.colors[logrus.InfoLevel].Printf(format, args...)
	cl.Logger.Infof(format, args...)
}

func (cl *ColorfulLogger) Warnf(format string, args ...interface{}) {
	cl.colors[logrus.WarnLevel].Printf(format, args...)
	cl.Logger.Warnf(format, args...)
}

func (cl *ColorfulLogger) Errorf(format string, args ...interface{}) {
	cl.colors[logrus.ErrorLevel].Printf(format, args...)
	cl.Logger.Errorf(format, args...)
}

func (cl *ColorfulLogger) Fatalf(format string, args ...interface{}) {
	cl.colors[logrus.FatalLevel].Printf(format, args...)
	cl.Logger.Fatalf(format, args...)
}

func (cl *ColorfulLogger) Panicf(format string, args ...interface{}) {
	cl.colors[logrus.PanicLevel].Printf(format, args...)
	cl.Logger.Panicf(format, args...)
}
