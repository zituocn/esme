package logx

import (
	"io"
	"time"
)

type Writer interface {
	WriteLog(time.Time, int, []byte) (int, error)
}

type writer struct {
	w io.Writer
}

func NewWriter(w io.Writer) Writer {
	return writer{w: w}
}

func (wr writer) WriteLog(t time.Time, level int, p []byte) (int, error) {
	return wr.w.Write(p)
}
