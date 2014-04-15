/*
 * Implementation  of the standard syslog client functionality
 *
 * This module is a part of the logger. It implements the subset of the
 * standard syslog client functionality (the subset I need...)
 */
package utils

import (
	"errors"
	"fmt"
	"net"
	"time"
)
/* Defines syslog severity */
type Severity int
const (
	Emergency Severity = iota
	Alert
	Critical
	Error
	Warning
	Notice
	Informational
	Debug
)
/* Usual string representation. */
func (s Severity) String() string {
	switch s {
	case Emergency:
		return "EMERGENCY"
	case Alert:
		return "ALERT"
	case Critical:
		return "CRITICAL"
	case Error:
		return "ERROR"
	case Warning:
		return "WARNING"
	case Notice:
		return "NOTICE"
	case Informational:
		return "INFO"
	case Debug:
		return "DEBUG"
	default:
		panic(errors.New("syslog: Invalid Severity values"))
	}
	return ""
}
/* Implementation of the syslog facility. */
type Facility int
const (
	FacKernel Facility = iota
	FacUser
	FacMail
	FacSystem
	FacSecurity4
	FacSyslogd
	FacLine
	FacNetwork
	FacUUCP
	FacClock9
	FacSecurity10
	FacFTP
	FacNTP
	FacLogAudit
	FacLogAlert
	FacClock15
	FacLocal0
	FacLocal1
	FacLocal2
	FacLocal3
	FacLocal4
	FacLocal5
	FacLocal6
	FacLocal7
)

const (
	// Define a standard syslog message timestamp format
	TimestampFmt = "Jan _2 15:04:05"

	// Standard UDP port for syslog is 514
	SyslogPort = 514
)

/* Implementation of the complete syslog messagea.  */
type SyslogMsg struct {
	Sev                 Severity
	Fac                 Facility
	timestamp, Hostname string
	Msg                 string
}
/* Returns syslog message's priority as a string. */
func (s *SyslogMsg) Priority() string {
	pri := int(s.Sev) + (8 * int(s.Fac))
	return fmt.Sprintf("<%d>", pri)
}

/* Returns syslog message's timestamp as a string. */
func (s *SyslogMsg) TimeStamp() string { return s.timestamp }

/* Set syslog message's timestamp. */
func (s *SyslogMsg) SetTimestamp(stamp time.Time) {
	s.timestamp = stamp.Format(TimestampFmt)
}

/* Set syslog message's timestamp. Timestamp is defined as a string. */
func (s *SyslogMsg) SSetTimestamp(stamp string) error {
	t, err := time.Parse(TimestampFmt, stamp)
	if err != nil { return err }
	s.SetTimestamp(t)
	return nil
}

/* Return the complete syslog message. */
func (s *SyslogMsg) Get() string {
	format := "%s%s %s %s"
	return fmt.Sprintf(format, s.Priority(), s.timestamp, s.Hostname, s.Msg)
}

/* Send a syslog message to given IP address. */
func (s *SyslogMsg) Send(ip string) error {
	var addr net.IP
	// local IP address overrides the Hostname field
	if ip != "" { s.Hostname = ip }
	addr = net.ParseIP(s.Hostname)
	// let's make an UDP connection and send the message
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{addr, SyslogPort, ""})
	if err != nil { return err }
	defer conn.Close()
	fmt.Fprintf(conn, s.Get())
	return nil
}

/* Create new syslog message with default fields */
func NewSyslogMsg() *SyslogMsg {
	return &SyslogMsg{Informational, FacLocal0, "", "", ""}
}
