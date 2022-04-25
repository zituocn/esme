package logx

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	LogDate = 1 << iota
	LogTime
	LogMicroSeconds
	LogLongFile
	LogShortFile
	LogUTC
	LogModule
	LogLevel

	StdFlags = LogDate | LogMicroSeconds | LogShortFile | LogLevel
)

const (
	TEST = iota
	DEBUG
	INFO
	NOTICE
	WARN
	ERROR
	PANIC
	FATAL
)

var (
	levels = []string{
		"[T]",
		"[D]",
		"[I]",
		"[N]",
		"[W]",
		"[E]",
		"[P]",
		"[F]",
	}
)

type Logger struct {
	flag      int
	level     int
	out       Writer
	callDepth int
	prefix    string
	Pool      *sync.Pool
}

func NewLogger(w Writer, flag, level int) *Logger {
	return &Logger{
		flag:      flag,
		level:     level,
		out:       w,
		callDepth: 2,
		Pool: &sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(nil)
			},
		},
	}
}

func (log *Logger) Info(format string, v ...interface{}) {
	if INFO < log.level {
		return
	}
	log.output(INFO, fmt.Sprintf(format, v...))
}

func (log *Logger) Debug(format string, v ...interface{}) {
	if DEBUG < log.level {
		return
	}
	log.output(DEBUG, fmt.Sprintf(format, v...))
}

func (log *Logger) Notice(format string, v ...interface{}) {
	if NOTICE < log.level {
		return
	}
	log.output(NOTICE, fmt.Sprintf(format, v...))
}

func (log *Logger) Error(format string, v ...interface{}) {
	if ERROR < log.level {
		return
	}
	log.output(ERROR, fmt.Sprintf(format, v...))
}

func (log *Logger) Warn(format string, v ...interface{}) {
	if WARN < log.level {
		return
	}
	log.output(WARN, fmt.Sprintf(format, v...))
}

func (log *Logger) Panic(format string, v ...interface{}) {
	if PANIC < log.level {
		return
	}
	s := fmt.Sprintf(format, v...)
	log.output(PANIC, s)
	panic(s)
}

func (log *Logger) Fatal(format string, v ...interface{}) {
	if FATAL < log.level {
		return
	}
	log.output(FATAL, fmt.Sprintf(format, v...))
	os.Exit(-1)
}

/*
get set
*/

func (log *Logger) GetCallDepth() int {
	return log.callDepth
}

func (log *Logger) SetCalldepth(depth int) {
	log.callDepth = depth
}

func (log *Logger) SetFlags(flag int) {
	log.flag = flag
}

func (log *Logger) SetPrefix(prefix string) {
	log.prefix = prefix
}

func (log *Logger) SetOutput(w Writer, prefix string) {
	log.out = w
	log.prefix = prefix
}

/*
private
*/

func (log *Logger) output(level int, msg string) {
	var (
		now  = time.Now()
		file string
		line int
	)

	if log.flag&(LogShortFile|LogLongFile) != 0 {
		ok := false
		_, file, line, ok = runtime.Caller(log.callDepth)
		if !ok {
			file = "???"
			line = 0
		}
	}

	buffer := log.Pool.Get().(*bytes.Buffer)
	buffer.Reset()

	log.writeHeader(buffer, now, file, line, level)

	buffer.WriteString(msg)

	if len(msg) > 0 && msg[len(msg)-1] != '\n' {
		buffer.WriteByte('\n')
	}

	log.out.WriteLog(now, level, buffer.Bytes())

	log.Pool.Put(buffer)
}

func (log *Logger) writeHeader(buffer *bytes.Buffer, now time.Time, file string, line int, level int) {
	// prefix
	if log.prefix != "" {
		buffer.WriteByte('[')
		buffer.WriteString(log.prefix)
		buffer.WriteByte(']')
		buffer.WriteByte(' ')
	}

	// datetime
	if log.flag&(LogDate|LogTime|LogMicroSeconds) != 0 {

		if log.flag&LogDate != 0 {
			year, month, day := now.Date()
			formatWrite(buffer, year, 4)
			buffer.WriteByte('/')

			formatWrite(buffer, int(month), 2)
			buffer.WriteByte('/')

			formatWrite(buffer, day, 2)
			buffer.WriteByte(' ')
		}

		if log.flag&(LogTime|LogMicroSeconds) != 0 {
			hour, min, sec := now.Clock()
			formatWrite(buffer, hour, 2)
			buffer.WriteByte(':')

			formatWrite(buffer, min, 2)
			buffer.WriteByte(':')

			formatWrite(buffer, sec, 2)

			if log.flag&LogMicroSeconds != 0 {
				buffer.WriteByte('.')
				formatWrite(buffer, now.Nanosecond()/1e6, 3)
			}
			buffer.WriteByte(' ')
		}
	}

	// level
	if log.flag&LogLevel != 0 {
		buffer.WriteString(levels[level])
		buffer.WriteByte(' ')
	}

	// package
	if log.flag&LogModule != 0 {
		buffer.WriteByte('[')
		buffer.WriteString(moduleOf(file))
		buffer.WriteByte(']')
		buffer.WriteByte(' ')
	}

	// filename and line
	if log.flag&(LogShortFile|LogLongFile) != 0 {
		if log.flag&LogShortFile != 0 {
			i := strings.LastIndex(file, "/")
			file = file[i+1:]
		}
		buffer.WriteString(file)
		buffer.WriteByte(':')
		formatWrite(buffer, line, -1)
		buffer.WriteByte(':')
		buffer.WriteByte(' ')
	}
}

func formatWrite(buffer *bytes.Buffer, i int, wid int) {
	var u = uint(i)
	if u == 0 && wid <= 1 {
		buffer.WriteByte('0')
		return
	}
	var b [32]byte
	bp := len(b)
	for ; u > 0 || wid > 0; u /= 10 {
		bp--
		wid--
		b[bp] = byte(u%10) + '0'
	}

	for bp < len(b) {
		buffer.WriteByte(b[bp])
		bp++
	}
}

func moduleOf(file string) string {
	pos := strings.LastIndex(file, "/")
	if pos != -1 {
		pos1 := strings.LastIndex(file[:pos], "/src/")
		if pos1 != -1 {
			return file[pos1+5 : pos]
		}
	}

	return "UNKNOWN"
}
