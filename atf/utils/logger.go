/*
 * logger.go -  custom logger implementation
 *
 * Default log module is simply not enough
 *
 * History:
 *  0.1.0   Jul11   MR  The initial version
 */

package utils

import (
	"fmt"
	"os"
	"strings"
	"time"
)

/************************** LogLevel ***********************************/
/*
 * LogLevel - an enum defining log levels
 */
type LogLevel Severity

const (
	EmergencyLogLevel LogLevel = iota
	AlertLogLevel
	CriticalLogLevel
	ErrorLogLevel
	WarningLogLevel
	NoticeLogLevel
	InfoLogLevel
	DebugLogLevel
	UnknownLogLevel
)

/*
 * LogLevel.String - a method returning the string representation of the
 *                   LogLevel value
 */
func (ll LogLevel) String() (s string) {
	switch ll {
	case EmergencyLogLevel:
		s = "EMERGENCY"
	case AlertLogLevel:
		s = "ALERT"
	case CriticalLogLevel:
		s = "CRITICAL"
	case ErrorLogLevel:
		s = "ERROR"
	case WarningLogLevel:
		s = "WARNING"
	case NoticeLogLevel:
		s = "NOTICE"
	case InfoLogLevel:
		s = "INFO"
	case DebugLogLevel:
		s = "DEBUG"
	default:
		panic("Unknown Log Level")
	}
	return s
}

/*
 * LogLevelFromString - converts log level given as string into proper LogLevel
 *                      value
 *
 * If invalid string is given, function returns 'UnknownLogLevel' value.
 */
func LogLevelFromString(lvl string) LogLevel {
	loglvl := UnknownLogLevel
	switch strings.ToUpper(lvl) {
	case "EMERGENCY":
		loglvl = EmergencyLogLevel
	case "ALERT":
		loglvl = AlertLogLevel
	case "CRITICAL":
		loglvl = CriticalLogLevel
	case "ERROR":
		loglvl = ErrorLogLevel
	case "WARNING":
		loglvl = WarningLogLevel
	case "NOTICE":
		loglvl = NoticeLogLevel
	case "INFO":
		loglvl = InfoLogLevel
	case "DEBUG":
		loglvl = DebugLogLevel
	}
	return loglvl
}

/************************** Logger ***********************************/
/*
 * Logger - an interface defining methods for various log handlers
 */
type Logger interface {
	Emergency(string)
	Alert(string)
	Critical(string)
	Error(string)
	Warning(string)
	Notice(string)
	Info(string)
	Debug(string)
	String() string
	Close()
}

/************************** logHandler ***********************************/
/*
 * logHandler - private struct that defines log level and format
 */
type logHandler struct {
	level  LogLevel /* log level set for this LogHandler */
	format string   /* a formatter for this LogHandler */
}

func (l *logHandler) Level() LogLevel { return l.level }

func (l *logHandler) SetLevel(lvl LogLevel) { l.level = lvl }

func (l *logHandler) Format() string { return l.format }

func (l *logHandler) SetFormat(fmt string) { l.format = fmt }

func newLogHandler(fmt string, lvl LogLevel) *logHandler {
	return &logHandler{lvl, fmt}
}

/************************** Log ***********************************/
/*
 * Log - a slice of different Loggers that can be added at will
 */
type Log struct {
	Handlers []Logger
}

func (l *Log) String() string {
	s := ""
	for _, h := range l.Handlers {
		if h != nil {
			s += fmt.Sprint(h.String())
		}
	}
	return s
}

func (l *Log) findEmpty() int {
	for ix, h := range l.Handlers {
		if h == nil {
			return ix
		}
	}
	return -1
}

func (l *Log) AddHandler(log Logger) []Logger {
	// check length and capacity; resize if needed...
	length := len(l.Handlers)
	c := cap(l.Handlers)
	if length+1 > c {
		newlst := make([]Logger, 0, 2*c)
		copy(newlst, l.Handlers)
		l.Handlers = newlst
	}
	l.Handlers = l.Handlers[0 : length+1]
	// now find empty index and insert Handler
	ix := l.findEmpty()
	if ix > -1 {
		l.Handlers[ix] = log
	}
	return l.Handlers
}

/*
 * Log - a generic Log method
 * 
 * Calls all needed log handlers and logs the given message with given level.
 * If an unknown log level is specified, do nothing.
 */
func (l *Log) Log(level LogLevel, msg string) {
	for _, h := range l.Handlers {
		switch level {
		case EmergencyLogLevel:
			h.Emergency(msg)
		case AlertLogLevel:
			h.Alert(msg)
		case CriticalLogLevel:
			h.Critical(msg)
		case ErrorLogLevel:
			h.Error(msg)
		case WarningLogLevel:
			h.Warning(msg)
		case NoticeLogLevel:
			h.Notice(msg)
		case InfoLogLevel:
			h.Info(msg)
		case DebugLogLevel:
			h.Debug(msg)
		}
	}
}

/*
 * LogS - a string version of the Log() method (see above)
 *
 * Calls all needed log handlers and logs the given message with given level.
 * Level is specified as string. 
 * If an unknown log level is specified, do nothing.
 */
func (l *Log) LogS(level string, msg string) {
	// get a Loglevel value from given string
	ll := LogLevelFromString(level)

	// check that a valid LogLevel value has been received and log the message;
	// if invalid log level, do nothing
	if ll != UnknownLogLevel {
		l.Log(ll, msg)
	}
}

func (l *Log) Debug(msg string) {
	for _, h := range l.Handlers {
		h.Debug(msg)
	}
}

func (l *Log) Info(msg string) {
	for _, h := range l.Handlers {
		h.Info(msg)
	}
}

func (l *Log) Notice(msg string) {
	for _, h := range l.Handlers {
		h.Notice(msg)
	}
}

func (l *Log) Warning(msg string) {
	for _, h := range l.Handlers {
		h.Warning(msg)
	}
}

func (l *Log) Error(msg string) {
	for _, h := range l.Handlers {
		h.Error(msg)
	}
}

func (l *Log) Critical(msg string) {
	for _, h := range l.Handlers {
		h.Critical(msg)
	}
}

func (l *Log) Alert(msg string) {
	for _, h := range l.Handlers {
		h.Alert(msg)
	}
}

func (l *Log) Emergency(msg string) {
	for _, h := range l.Handlers {
		h.Emergency(msg)
	}
}

func (l *Log) Len() int { return len(l.Handlers) }

func (l *Log) Close() {
	for _, h := range l.Handlers {
		h.Close()
	}
}

const logLength int = 5

func NewLog(num int) *Log {
	if num == 0 {
		num = logLength
	}
	return &Log{make([]Logger, 0, logLength)}
}

/************************** Formatter  ***********************************/
/*
 * Formatter - an interface defining the generic formatter
 */
type Formatter interface {
	Format(string)
}

/************************** FileHandler ***********************************/
/*
 * FileHandler
 */
type FileHandler struct {
	*logHandler
	file *os.File
}

/*
 * FileHandler.log - creates new stream handler
 */
func (f *FileHandler) log(level LogLevel, msg string) {
	if f.Level() >= level {
		fmt.Fprintf(f.file, f.Format(), Now(), level, msg)
	}
}

func (f *FileHandler) Debug(msg string) {
	f.log(DebugLogLevel, msg)
}

func (f *FileHandler) Info(msg string) {
	f.log(InfoLogLevel, msg)
}

func (f *FileHandler) Notice(msg string) {
	f.log(NoticeLogLevel, msg)
}

func (f *FileHandler) Warning(msg string) {
	f.log(WarningLogLevel, msg)
}

func (f *FileHandler) Error(msg string) {
	f.log(ErrorLogLevel, msg)
}

func (f *FileHandler) Critical(msg string) {
	f.log(CriticalLogLevel, msg)
}

func (f *FileHandler) Alert(msg string) {
	f.log(AlertLogLevel, msg)
}

func (f *FileHandler) Emergency(msg string) {
	f.log(EmergencyLogLevel, msg)
}

func (f *FileHandler) Close() {
	if f.file != nil {
		f.file.Close()
	}
}

func (f *FileHandler) String() string {
	return fmt.Sprintf("  FileHandler: fmt=%q, lvl=%-10s, fd=%d\n",
		f.Format(), f.Level(), f.file.Fd())
}

/*
 * NewFileHandler - creates new file handler
 */
func NewFileHandler(filename string,
	fmt string, lvl LogLevel) (*FileHandler, error) {
	// open log file
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0755)
	return &FileHandler{newLogHandler(fmt, lvl), f}, err
}

/************************** StreamHandler ***********************************/
/*
 * StreamHandler
 */
type StreamHandler FileHandler

/*
 * StreamHandler.log - creates new stream handler
 */
func (s *StreamHandler) log(level LogLevel, msg string) {
	if s.Level() >= level {
		fmt.Printf(s.Format(), Now(), level, msg)
	}
}

func (s *StreamHandler) Debug(msg string) {
	s.log(DebugLogLevel, msg)
}

func (s *StreamHandler) Info(msg string) {
	s.log(InfoLogLevel, msg)
}

func (s *StreamHandler) Notice(msg string) {
	s.log(NoticeLogLevel, msg)
}

func (s *StreamHandler) Warning(msg string) {
	s.log(WarningLogLevel, msg)
}

func (s *StreamHandler) Error(msg string) {
	s.log(ErrorLogLevel, msg)
}

func (s *StreamHandler) Critical(msg string) {
	s.log(CriticalLogLevel, msg)
}

func (s *StreamHandler) Alert(msg string) {
	s.log(AlertLogLevel, msg)
}

func (s *StreamHandler) Emergency(msg string) {
	s.log(EmergencyLogLevel, msg)
}

func (s *StreamHandler) Close() {
	// an empty implementation to satisfy the Logger interface
}

func (s *StreamHandler) String() string {
	return fmt.Sprintf("StreamHandler: fmt=%q, lvl=%-10s\n",
		s.Format(), s.Level())
}

/*
 * NewStreamHandler - creates new stream handler
 */
func NewStreamHandler(fmt string, lvl LogLevel) *StreamHandler {
	return &StreamHandler{newLogHandler(fmt, lvl), os.Stdout}
}

/************************** SyslogHandler ***********************************/
/*
 * SyslogHandler - 
 */
type SyslogHandler struct {
	*logHandler
	IP string
	*SyslogMsg
}

/*
 * SyslogHandler.log - sends a message to the wire using UDP port 514
 */
func (s *SyslogHandler) log(level LogLevel, msg string) error {
	if s.Level() >= level {
		s.Fac = FacLocal0
		s.Sev = Severity(level)
		s.Msg = fmt.Sprintf("%s %s", level.String(), msg)
		t := time.Now()
		s.SetTimestamp(t)
		err := s.Send(s.IP)
		if err != nil {
			panic(err)
		}
	}
	return nil
}

func (s *SyslogHandler) Debug(msg string) {
	s.log(DebugLogLevel, msg)
}

func (s *SyslogHandler) Info(msg string) {
	s.log(InfoLogLevel, msg)
}

func (s *SyslogHandler) Notice(msg string) {
	s.log(NoticeLogLevel, msg)
}

func (s *SyslogHandler) Warning(msg string) {
	s.log(WarningLogLevel, msg)
}

func (s *SyslogHandler) Error(msg string) {
	s.log(ErrorLogLevel, msg)
}

func (s *SyslogHandler) Critical(msg string) {
	s.log(CriticalLogLevel, msg)
}

func (s *SyslogHandler) Alert(msg string) {
	s.log(AlertLogLevel, msg)
}

func (s *SyslogHandler) Emergency(msg string) {
	s.log(EmergencyLogLevel, msg)
}

func (s *SyslogHandler) Close() {
	// an empty implementation to satisfy the Logger interface
}

func (s *SyslogHandler) String() string {
	return fmt.Sprintf("SyslogHandler: fmt=%q, lvl=%-10s, Server=%q %s\n",
		s.Format(), s.Level(), s.IP)
}

/*
 * NewSyslogHandler - creates new syslog handler
 */
func NewSyslogHandler(ip string, fmt string, lvl LogLevel) *SyslogHandler {
	//
	return &SyslogHandler{newLogHandler(fmt, lvl), ip, NewSyslogMsg()}
}
