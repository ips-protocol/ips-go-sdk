package log

//import (
//	"bytes"
//	"fmt"
//	"io"
//	"os"
//	"runtime"
//	"strings"
//	"sync"
//	"time"
//)
//
//const (
//	Ldate = 1 << iota
//	Ltime
//	Lmicroseconds
//	Llongfile
//	Lshortfile
//	Lmodule
//	Llevel
//	LstdFlags = Ldate | Ltime | Lmicroseconds
//	Ldefault  = Lmodule | Llevel | Lshortfile | LstdFlags
//) // [prefix][time][level][module][shortfile|longfile]
//
//const (
//	Ldebug = iota
//	Linfo
//	Lwarn
//	Lerror
//	Lpanic
//	Lfatal
//)
//
//var levels = []string{
//	"[DEBUG]",
//	"[INFO]",
//	"[WARN]",
//	"[ERROR]",
//	"[PANIC]",
//	"[FATAL]",
//}
//
//type Logger struct {
//	mu         sync.Mutex
//	prefix     string
//	flag       int
//	Level      int
//	out        io.Writer
//	buf        bytes.Buffer
//	levelStats [6]int64
//}
//
//func New(out io.Writer, prefix string, flag int) *Logger {
//	return &Logger{out: out, prefix: prefix, Level: 1, flag: flag}
//}
//
//var Std = New(os.Stderr, "", Ldefault)
//
//func itoa(buf *bytes.Buffer, i int, wid int) {
//	var u uint = uint(i)
//	if u == 0 && wid <= 1 {
//		buf.WriteByte('0')
//		return
//	}
//
//	// Assemble decimal in reverse order.
//	var b [32]byte
//	bp := len(b)
//	for ; u > 0 || wid > 0; u /= 10 {
//		bp--
//		wid--
//		b[bp] = byte(u%10) + '0'
//	}
//
//	// avoid slicing b to avoid an allocation.
//	for bp < len(b) {
//		buf.WriteByte(b[bp])
//		bp++
//	}
//}
//
//func shortFile(file string, flag int) string {
//	sep := "/"
//	if (flag & Lmodule) != 0 {
//		sep = "/src/"
//	}
//	pos := strings.LastIndex(file, sep)
//	if pos != -1 {
//		return file[pos+5:]
//	}
//	return file
//}
//
//func (l *Logger) formatHeader(buf *bytes.Buffer, t time.Time, file string, line int, lvl int) {
//	if l.prefix != "" {
//		buf.WriteString(l.prefix)
//	}
//	if l.flag&(Ldate|Ltime|Lmicroseconds) != 0 {
//		if l.flag&Ldate != 0 {
//			year, month, day := t.Date()
//			itoa(buf, year, 4)
//			buf.WriteByte('/')
//			itoa(buf, int(month), 2)
//			buf.WriteByte('/')
//			itoa(buf, day, 2)
//			buf.WriteByte(' ')
//		}
//		if l.flag&(Ltime|Lmicroseconds) != 0 {
//			hour, min, sec := t.Clock()
//			itoa(buf, hour, 2)
//			buf.WriteByte(':')
//			itoa(buf, min, 2)
//			buf.WriteByte(':')
//			itoa(buf, sec, 2)
//			if l.flag&Lmicroseconds != 0 {
//				buf.WriteByte('.')
//				itoa(buf, t.Nanosecond()/1e3, 6)
//			}
//			buf.WriteByte(' ')
//		}
//	}
//	if l.flag&Llevel != 0 {
//		buf.WriteString(levels[lvl])
//	}
//	if l.flag&(Lshortfile|Llongfile) != 0 {
//		if l.flag&Lshortfile != 0 {
//			file = shortFile(file, l.flag)
//		}
//		buf.WriteByte(' ')
//		buf.WriteString(file)
//		buf.WriteByte(':')
//		itoa(buf, line, -1)
//		buf.WriteString(": ")
//	}
//}
//
//func (l *Logger) Output(lvl int, calldepth int, s string) error {
//	if lvl < l.Level {
//		return nil
//	}
//	now := time.Now() // get this early.
//	var file string
//	var line int
//	l.mu.Lock()
//	defer l.mu.Unlock()
//	if l.flag&(Lshortfile|Llongfile|Lmodule) != 0 {
//		// release lock while getting caller info - it's expensive.
//		l.mu.Unlock()
//		var ok bool
//		_, file, line, ok = runtime.Caller(calldepth)
//		if !ok {
//			file = "???"
//			line = 0
//		}
//		l.mu.Lock()
//	}
//	l.levelStats[lvl]++
//	l.buf.Reset()
//	l.formatHeader(&l.buf, now, file, line, lvl)
//	l.buf.WriteString(s)
//	if len(s) > 0 && s[len(s)-1] != '\n' {
//		l.buf.WriteByte('\n')
//	}
//	_, err := l.out.Write(l.buf.Bytes())
//	return err
//}
//
//func (l *Logger) Printf(format string, v ...interface{}) {
//	l.Output(Linfo, 2, fmt.Sprintf(format, v...))
//}
//
//func (l *Logger) Print(v ...interface{}) { l.Output(Linfo, 2, fmt.Sprint(v...)) }
//
//func (l *Logger) Println(v ...interface{}) { l.Output(Linfo, 2, fmt.Sprintln(v...)) }
//
//func (l *Logger) Debugf(format string, v ...interface{}) {
//	if Ldebug < l.Level {
//		return
//	}
//	l.Output(Ldebug, 2, fmt.Sprintf(format, v...))
//}
//
//func (l *Logger) Debug(v ...interface{}) {
//	if Ldebug < l.Level {
//		return
//	}
//	l.Output(Ldebug, 2, fmt.Sprintln(v...))
//}
//
//func (l *Logger) Infof(format string, v ...interface{}) {
//	if Linfo < l.Level {
//		return
//	}
//	l.Output(Linfo, 2, fmt.Sprintf(format, v...))
//}
//
//func (l *Logger) Info(v ...interface{}) {
//	if Linfo < l.Level {
//		return
//	}
//	l.Output(Linfo, 2, fmt.Sprintln(v...))
//}
//
//func (l *Logger) Warnf(format string, v ...interface{}) {
//	l.Output(Lwarn, 2, fmt.Sprintf(format, v...))
//}
//
//func (l *Logger) Warn(v ...interface{}) { l.Output("", Lwarn, 2, fmt.Sprintln(v...)) }
//
//func (l *Logger) Errorf(format string, v ...interface{}) {
//	l.Output(Lerror, 2, fmt.Sprintf(format, v...))
//}
//
//func (l *Logger) Error(v ...interface{}) { l.Output("", Lerror, 2, fmt.Sprintln(v...)) }
//
//func (l *Logger) Fatal(v ...interface{}) {
//	l.Output(Lfatal, 2, fmt.Sprint(v...))
//	os.Exit(1)
//}
//
//func (l *Logger) Fatalf(format string, v ...interface{}) {
//	l.Output(Lfatal, 2, fmt.Sprintf(format, v...))
//	os.Exit(1)
//}
//
//func (l *Logger) Fatalln(v ...interface{}) {
//	l.Output(Lfatal, 2, fmt.Sprintln(v...))
//	os.Exit(1)
//}
//
//func (l *Logger) Panic(v ...interface{}) {
//	s := fmt.Sprint(v...)
//	l.Output(Lpanic, 2, s)
//	panic(s)
//}
//
//func (l *Logger) Panicf(format string, v ...interface{}) {
//	s := fmt.Sprintf(format, v...)
//	l.Output(Lpanic, 2, s)
//	panic(s)
//}
//
//func (l *Logger) Stack(v ...interface{}) {
//	s := fmt.Sprint(v...)
//	s += "\n"
//	buf := make([]byte, 1024*1024)
//	n := runtime.Stack(buf, true)
//	s += string(buf[:n])
//	s += "\n"
//	l.Output(Lerror, 2, s)
//}
//
//func (l *Logger) SingleStack(v ...interface{}) {
//	s := fmt.Sprint(v...)
//	s += "\n"
//	buf := make([]byte, 1024*1024)
//	n := runtime.Stack(buf, false)
//	s += string(buf[:n])
//	s += "\n"
//	l.Output(Lerror, 2, s)
//}
//
//func (l *Logger) Stat() (stats []int64) {
//	l.mu.Lock()
//	v := l.levelStats
//	l.mu.Unlock()
//	return v[:]
//}
//
//func (l *Logger) Flags() int {
//	l.mu.Lock()
//	defer l.mu.Unlock()
//	return l.flag
//}
//
//func (l *Logger) SetFlags(flag int) {
//	l.mu.Lock()
//	defer l.mu.Unlock()
//	l.flag = flag
//}
//
//func (l *Logger) Prefix() string {
//	l.mu.Lock()
//	defer l.mu.Unlock()
//	return l.prefix
//}
//
//func (l *Logger) SetPrefix(prefix string) {
//	l.mu.Lock()
//	defer l.mu.Unlock()
//	l.prefix = prefix
//}
//
//func (l *Logger) SetOutputLevel(lvl int) {
//	l.mu.Lock()
//	defer l.mu.Unlock()
//	l.Level = lvl
//}
//
//func SetOutput(w io.Writer) {
//	Std.mu.Lock()
//	defer Std.mu.Unlock()
//	Std.out = w
//}
//
//func Flags() int {
//	return Std.Flags()
//}
//
//func SetFlags(flag int) {
//	Std.SetFlags(flag)
//}
//
//func Prefix() string {
//	return Std.Prefix()
//}
//
//func SetPrefix(prefix string) {
//	Std.SetPrefix(prefix)
//}
//
//func SetOutputLevel(lvl int) {
//	Std.SetOutputLevel(lvl)
//}
//
//func GetOutputLevel() int {
//	return Std.Level
//}
//
//func Print(v ...interface{}) {
//	Std.Output(Linfo, 2, fmt.Sprint(v...))
//}
//
//func Printf(format string, v ...interface{}) {
//	Std.Output(Linfo, 2, fmt.Sprintf(format, v...))
//}
//
//func Println(v ...interface{}) {
//	Std.Output(Linfo, 2, fmt.Sprintln(v...))
//}
//
//func Debugf(format string, v ...interface{}) {
//	if Ldebug < Std.Level {
//		return
//	}
//	Std.Output(Ldebug, 2, fmt.Sprintf(format, v...))
//}
//
//func Debug(v ...interface{}) {
//	if Ldebug < Std.Level {
//		return
//	}
//	Std.Output(Ldebug, 2, fmt.Sprintln(v...))
//}
//
//func Infof(format string, v ...interface{}) {
//	if Linfo < Std.Level {
//		return
//	}
//	Std.Output(Linfo, 2, fmt.Sprintf(format, v...))
//}
//
//func Info(v ...interface{}) {
//	if Linfo < Std.Level {
//		return
//	}
//	Std.Output(Linfo, 2, fmt.Sprintln(v...))
//}
//
//func Warnf(format string, v ...interface{}) {
//	Std.Output(Lwarn, 2, fmt.Sprintf(format, v...))
//}
//
//func Warn(v ...interface{}) { Std.Output(Lwarn, 2, fmt.Sprintln(v...)) }
//
//func Errorf(format string, v ...interface{}) {
//	Std.Output(Lerror, 2, fmt.Sprintf(format, v...))
//}
//
//func Error(v ...interface{}) { Std.Output(Lerror, 2, fmt.Sprintln(v...)) }
//
//func Fatal(v ...interface{}) {
//	Std.Output(Lfatal, 2, fmt.Sprint(v...))
//	os.Exit(1)
//}
//
//func Fatalf(format string, v ...interface{}) {
//	Std.Output(Lfatal, 2, fmt.Sprintf(format, v...))
//	os.Exit(1)
//}
//
//func Panic(v ...interface{}) {
//	s := fmt.Sprint(v...)
//	Std.Output(Lpanic, 2, s)
//	panic(s)
//}
//
//func Panicf(format string, v ...interface{}) {
//	s := fmt.Sprintf(format, v...)
//	Std.Output(Lpanic, 2, s)
//	panic(s)
//}
//
//func Panicln(v ...interface{}) {
//	s := fmt.Sprintln(v...)
//	Std.Output(Lpanic, 2, s)
//	panic(s)
//}
//
//func Stack(v ...interface{}) {
//	Std.Stack(v...)
//}
//
//func SingleStack(v ...interface{}) {
//	Std.SingleStack(v...)
//}
