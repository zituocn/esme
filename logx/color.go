package logx

import (
	"sync"
	"time"
)

const (
	endColor = "\033[0m"
)

var (
	red     = "\033[1;31m" //红色
	green   = "\033[1;32m" //绿色
	yellow  = "\033[1;33m" //黄色
	purple  = "\033[1;34m" //紫色
	magenta = "\033[1;35m" //品红
	teal    = "\033[0;97m" //青色
	white   = "\033[1;37m" //白色
)

var (
	logColor = []string{
		info:     white,
		debug:    purple,
		notice:   green,
		warn:     yellow,
		logError: red,
		logPanic: magenta,
		fatal:    magenta,
	}
)

type colorWriter struct {
	mu     *sync.Mutex
	b      []byte
	writer Writer
}

// WriteLog colorWrite WriteLog
func (w *colorWriter) WriteLog(now time.Time, level int, b []byte) (int, error) {
	w.b = w.b[:0]
	w.b = append(w.b, logColor[level]...)
	w.b = append(w.b, b...)
	w.b = append(w.b, endColor...)
	return w.writer.WriteLog(now, level, w.b)
}

// WithColor use colorWriter
func WithColor(w Writer) Writer {
	return &colorWriter{writer: w}
}
