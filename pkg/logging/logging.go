package logging

import (
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

// ColorfulLogger is a custom logger implementation based on logrus.Logger.
// It adds colorized output to the log messages based on the log level.
// ColorfulLogger provides methods for setting log level, output destination, and logging messages at different levels.
//
// Usage Example:
//
//	logger := NewColorfulLogger()
//	logger.SetLevel(logrus.DebugLevel)
//	logger.SetOutput(os.Stdout)
//	logger.Debug("Debug message")
//	logger.Info("Info message")
//	logger.Warn("Warn message")
//	logger.Error("Error message")
//	logger.Fatal("Fatal message")
//	logger.Panic("Panic message")
//	logger.Debugf("Debug message: %s", "value")
//	logger.Infof("Info message: %s", "value")
//	logger.Warnf("Warn message: %s", "value")
//	logger.Errorf("Error message: %s", "value")
//	logger.Fatalf("Fatal message: %s", "value")
//	logger.Panicf("Panic message: %s", "value")
type ColorfulLogger struct {
	*logrus.Logger
	colors map[logrus.Level]*color.Color
}

// NewColorfulLogger creates a new instance of ColorfulLogger with default settings.
// It configures the logger with a map of log levels to colors and sets the formatter
// to use a text formatter with full timestamps. The logger's output is set to
// standard output and the log level is set to InfoLevel.
//
// Returns a pointer to the newly created ColorfulLogger.
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

// SetLevel sets the log level of the ColorfulLogger.
// It updates the log level of the underlying Logger instance with the provided "level" parameter.
func (cl *ColorfulLogger) SetLevel(level logrus.Level) {
	cl.Logger.SetLevel(level)
}

// SetOutput sets the output destination for the logger.
// It updates the output destination of the underlying Logger instance
// with the provided `output` parameter.
func (cl *ColorfulLogger) SetOutput(output io.Writer) {
	cl.Logger.SetOutput(output)
}

// Debug logs a message at the debug level.
// It prints the message using the color associated with the debug level,
// and then delegates to the underlying Logger to log the message as well.
func (cl *ColorfulLogger) Debug(args ...interface{}) {
	cl.colors[logrus.DebugLevel].Println(args...)
	cl.Logger.Debug(args...)
}

// Info logs a message at the info level.
// It prints the message using the color associated with the info level,
// and then delegates to the underlying Logger to log the message as well.
func (cl *ColorfulLogger) Info(args ...interface{}) {
	_, err := cl.colors[logrus.InfoLevel].Println(args...)
	if err != nil {
		return
	}
	cl.Logger.Info(args...)
}

// Warn logs a message at Warn level, both with colored output and the regular logger.
func (cl *ColorfulLogger) Warn(args ...interface{}) {
	_, err := cl.colors[logrus.WarnLevel].Println(args...)
	if err != nil {
		return
	}
	cl.Logger.Warn(args...)
}

// Error logs a message at the Error level. It prints the message using the corresponding color
// set for the Error level and also logs it using the underlying Logger.
func (cl *ColorfulLogger) Error(args ...interface{}) {
	_, err := cl.colors[logrus.ErrorLevel].Println(args...)
	if err != nil {
		return
	}
	cl.Logger.Error(args...)
}

// Fatal logs a message with log level Fatal using the configured logger instance and the associated color.
// This method prints the message to the console and then exits the program.
// It takes a variadic parameter args which represents the message to be logged.
func (cl *ColorfulLogger) Fatal(args ...interface{}) {
	_, err := cl.colors[logrus.FatalLevel].Println(args...)
	if err != nil {
		return
	}
	cl.Logger.Fatal(args...)
}

// Panic logs a message at level Panic on the ColorfulLogger and the underlying Logger.
// It also prints the message in color using the configured color.
func (cl *ColorfulLogger) Panic(args ...interface{}) {
	_, err := cl.colors[logrus.PanicLevel].Println(args...)
	if err != nil {
		return
	}
	cl.Logger.Panic(args...)
}

// Debugf formats and prints a debug level log message using the provided format and arguments.
// It first prints the colored log message using the color associated with the debug level,
// and then it delegates to the underlying logger's Debugf method to print the log message without color.
// This method is a member of the ColorfulLogger struct.
func (cl *ColorfulLogger) Debugf(format string, args ...interface{}) {
	_, err := cl.colors[logrus.DebugLevel].Printf(format, args...)
	if err != nil {
		return
	}
	cl.Logger.Debugf(format, args...)
}

// Infof writes an informational log message with a format string and arguments.
// It prints the log message using the color associated with the InfoLevel and
// calls the Infof method of the embedded logrus.Logger, passing the format
// string and arguments.
func (cl *ColorfulLogger) Infof(format string, args ...interface{}) {
	_, err := cl.colors[logrus.InfoLevel].Printf(format, args...)
	if err != nil {
		return
	}
	cl.Logger.Infof(format, args...)
}

// Warnf logs a formatted warning message with the given format and arguments.
// It prints the message using the configured log level color and also logs it
// using the underlying logger with the warn level.
func (cl *ColorfulLogger) Warnf(format string, args ...interface{}) {
	_, err := cl.colors[logrus.WarnLevel].Printf(format, args...)
	if err != nil {
		return
	}
	cl.Logger.Warnf(format, args...)
}

// Errorf formats and prints an error level log message using the provided format and arguments.
// It first prints the colored log message using the color associated with the error level,
// and then it delegates to the underlying logger's Errorf method to print the log message without color.
func (cl *ColorfulLogger) Errorf(format string, args ...interface{}) {
	_, err := cl.colors[logrus.ErrorLevel].Printf(format, args...)
	if err != nil {
		return
	}
	cl.Logger.Errorf(format, args...)
}

// Fatalf logs a message with the specified format and arguments at the fatal log level.
// It prints the formatted message to the console with colors, if configured,
// and then calls the underlying Logger's Fatalf method to log the message without colors.
func (cl *ColorfulLogger) Fatalf(format string, args ...interface{}) {
	_, err := cl.colors[logrus.FatalLevel].Printf(format, args...)
	if err != nil {
		return
	}
	cl.Logger.Fatalf(format, args...)
}

// Panicf logs a message at the Panic level to the logger.
// It receives a format string and a list of arguments,
// which will be formatted according to the format string and
// then logged at the Panic level using the logger's Printf method.
func (cl *ColorfulLogger) Panicf(format string, args ...interface{}) {
	_, err := cl.colors[logrus.PanicLevel].Printf(format, args...)
	if err != nil {
		return
	}
	cl.Logger.Panicf(format, args...)
}
