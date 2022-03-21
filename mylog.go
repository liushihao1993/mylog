package mylog

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	levelInfo = level("INFO")
	levelWarn = level("WARN")
	levelErr  = level("ERROR")
)

// Ctx .
func Ctx(ctx context.Context) *entry {
	pc := make([]uintptr, 2)
	_ = runtime.Callers(2, pc)
	frame := runtime.CallersFrames(pc)
	var file, function string
	var line int
	if f, ok := frame.Next(); ok {
		function = f.Function[strings.LastIndex(f.Function, "/")+1:]
		file = filepath.Base(f.File)
		line = f.Line
	}
	fields := make([][2]interface{}, 0, len(logStd.ctxKeys)+2)

	for _, k := range logStd.ctxKeys {
		if v := ctx.Value(k); v != nil {
			fields = append(fields, [2]interface{}{k, v})
		}
	}
	return &entry{
		time:     time.Now().Format("2006-01-02 15:04:05.000"),
		file:     file,
		line:     line,
		function: function,
		fields:   fields,
	}
}

var outBufPool = sync.Pool{New: func() interface{} { return bytes.NewBuffer(make([]byte, 0, 4<<10)) }} // default 4 kb

type logger struct {
	once       sync.Once // to init only once
	infoWriter io.Writer
	warnWriter io.Writer
	errWriter  io.Writer
	ctxKeys    []string // 存放与ctx中的链路追踪信息
}
type entry struct {
	level    level
	time     string
	file     string
	line     int
	function string
	msg      string
	fields   [][2]interface{}
}

var logStd = logger{
	once:       sync.Once{},
	infoWriter: os.Stdout,
	warnWriter: os.Stdout,
	errWriter:  os.Stdout,
}

type level string

// WithField key-value.
func (e *entry) WithField(k string, v interface{}) *entry {
	e.fields = append(e.fields, [2]interface{}{k, v})
	return e
}

// WithFields key-value.
func (e *entry) WithFields(k1 string, v1 interface{}, k2 string, v2 interface{}, kvs ...interface{}) *entry {
	e.fields = append(e.fields, [2]interface{}{k1, v1}, [2]interface{}{k2, v2})
	if len(kvs) == 0 {
		return e
	}
	if len(kvs)%2 != 0 {
		for i := 0; i < len(kvs); i += 2 {
			if i+1 < len(kvs) {
				e.fields = append(e.fields, [2]interface{}{fmt.Sprint(kvs[i]), kvs[i+1]})
			} else {
				e.fields = append(e.fields, [2]interface{}{fmt.Sprint(kvs[i]), "??? key或value缺失"})
			}
		}
		return e
	}
	for i := 0; i < len(kvs); i += 2 {
		e.fields = append(e.fields, [2]interface{}{kvs[i], kvs[i+1]})
	}
	return e
}

func (e *entry) Info(a ...interface{}) {
	_ = e.outputLn(levelInfo, a...)
}
func (e *entry) Warn(a ...interface{}) {
	_ = e.outputLn(levelWarn, a...)
}
func (e *entry) Error(a ...interface{}) {
	_ = e.outputLn(levelErr, a...)
}

func (e *entry) Infof(format string, a ...interface{}) {
	_ = e.outputFln(levelInfo, format, a...)
}
func (e *entry) Warnf(format string, a ...interface{}) {
	_ = e.outputFln(levelWarn, format, a...)
}
func (e *entry) Errorf(format string, a ...interface{}) {
	_ = e.outputFln(levelErr, format, a...)
}

func (e *entry) outputLn(l level, a ...interface{}) error {
	return e.output(l, "", a...)
}
func (e *entry) outputFln(l level, format string, a ...interface{}) error {
	return e.output(l, format, a...)
}

func (e *entry) writeOut(writer io.Writer) (int, error) {
	outBuf := outBufPool.Get().(*bytes.Buffer)
	outBuf.Reset()
	defer func() { outBufPool.Put(outBuf) }()
	outBuf.WriteString(fmt.Sprintf("[%s] %s %s:%d %s %s", e.level, e.time, e.file, e.line, e.function, e.msg))
	if len(e.fields) == 0 {
		outBuf.WriteString("\n")
		return writer.Write(outBuf.Bytes())
	}
	outBuf.WriteString(" {")
	var key string
	for i := 0; i < len(e.fields); i++ {
		if i != 0 {
			outBuf.WriteString(",")
		}
		key, _ = e.fields[i][0].(string)
		outBuf.WriteString(fmt.Sprintf(`"%s":`, key))
		switch s := e.fields[i][1].(type) {
		default:
			outBuf.WriteString(fmt.Sprintf(`%+v`, s))
		case []byte:
			outBuf.WriteString(fmt.Sprintf(`%s`, s))
		case *[]byte:
			outBuf.Write(*s)
		case *string:
			outBuf.WriteString(*s)
		}
	}
	outBuf.WriteString("}\n")
	return writer.Write(outBuf.Bytes())
}
func (e *entry) output(l level, format string, a ...interface{}) error {
	var msg string
	if format == "" {
		msg = fmt.Sprintln(a...)
		msg = msg[:len(msg)-1] // 去除末尾的\n符号
	} else {
		msg = fmt.Sprintf(format, a...)
	}
	e.msg = msg
	e.level = l

	var outWriter io.Writer
	switch l {
	case levelInfo:
		outWriter = logStd.infoWriter
	case levelWarn:
		outWriter = logStd.warnWriter
	case levelErr:
		outWriter = logStd.errWriter
	default: // to protect writer is not nil if l none of the above.
		outWriter = logStd.infoWriter
	}
	_, err := e.writeOut(outWriter)
	return err
}

// Init set the Writer and Context Keys.
func Init(infoF, warnF, errF io.Writer, CtxKeys ...string) {
	logStd.once.Do(func() {
		logStd.infoWriter = infoF
		logStd.warnWriter = warnF
		logStd.errWriter = errF
		logStd.ctxKeys = CtxKeys
	})
}
