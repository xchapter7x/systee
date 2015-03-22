package systee

import (
	"errors"
	"fmt"
	"log/syslog"

	"github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/hooks/syslog"
)

const (
	LogForwardMessage = "forwarded by systee"
)

var (
	ErrSyslogConnectionFailure = errors.New("Unable to connect to syslog daemon")
)

type LogHandler interface {
	Handle(LogMsg)
}

func NewLogforward(host string, port, proto int) (lh LogHandler, err error) {
	logForward := &Logforward{
		host:  host,
		port:  port,
		proto: proto,
	}

	if err = logForward.Connect(); err == nil {
		lh = logForward
	}
	return
}

type Logforward struct {
	host  string
	port  int
	proto int
}

func (s *Logforward) getProto() (p string) {
	switch s.proto {
	case TCP:
		p = "tcp"

	case UDP:
		p = "udp"
	}
	return
}

func (s *Logforward) Connect() (err error) {
	url := fmt.Sprintf("%s:%d", s.host, s.port)

	if hook, e := logrus_syslog.NewSyslogHook(s.getProto(), url, syslog.LOG_INFO, ""); e == nil {
		logrus.AddHook(hook)

	} else {
		err = ErrSyslogConnectionFailure
	}
	return
}

func (s *Logforward) Handle(logMessage LogMsg) {
	logrus.WithFields(logrus.Fields(logMessage)).Info(LogForwardMessage)
}
