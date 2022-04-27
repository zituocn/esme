package logx

import (
	"bytes"
	"os"
)

var (
	std *Logger
)

func init() {
	std = NewLogger(WithColor(NewWriter(os.Stdout)), StdFlags, debug)
	std.SetCalldepth(std.GetCallDepth() + 1)
}

func Info(v ...interface{}) {
	std.Info(format(v), v...)
}

func Debug(v ...interface{}) {
	std.Debug(format(v), v...)
}

func Notice(v ...interface{}) {
	std.Notice(format(v), v...)
}

func Warn(v ...interface{}) {
	std.Warn(format(v), v...)
}

func Error(v ...interface{}) {
	std.Error(format(v), v...)
}

func Fatal(v ...interface{}) {
	std.Fatal(format(v), v...)
}

func Panic(v ...interface{}) {
	std.Panic(format(v), v...)
}

/*
format
*/

func Infof(format string, v ...interface{}) {
	std.Info(format, v...)
}

func Noticef(format string, v ...interface{}) {
	std.Notice(format, v...)
}

func Debugf(format string, v ...interface{}) {
	std.Debug(format, v...)
}

func Warnf(format string, v ...interface{}) {
	std.Warn(format, v...)
}

func Errorf(format string, v ...interface{}) {
	std.Error(format, v...)
}

func Panicf(format string, v ...interface{}) {
	std.Panic(format, v...)
}

func Fatalf(format string, v ...interface{}) {
	std.Fatal(format, v...)
}

func format(v ...interface{}) string {
	length := len(v)
	buffer := &bytes.Buffer{}
	for i := 0; i < length; i++ {
		buffer.WriteString("%v ")
	}
	return buffer.String()
}
