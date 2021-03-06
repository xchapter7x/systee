package systee

import (
	"fmt"

	gsyslog "github.com/mcuadros/go-syslog"
	"gopkg.in/mcuadros/go-syslog.v2/format"
)

const (
	UDP = iota
	TCP
	TCPUDP
	RFC5424
	RFC3164
	RFC6587
)

func NewListener(host string, port, proto, format int) *Listener {
	return &Listener{
		host:   host,
		port:   port,
		proto:  proto,
		format: format,
	}
}

type Listener struct {
	port            int
	host            string
	proto           int
	format          int
	server          syslogServer
	logpartsChannel gsyslog.LogPartsChannel
}

type syslogServer interface {
	SetFormat(format.Format)
	ListenUDP(string) error
	ListenTCP(string) error
	Kill() error
}

type LogMsg map[string]interface{}

func (s *Listener) AddHandler(h func(LogMsg)) {

	go func() {
		for logParts := range s.logpartsChannel {
			h(LogMsg(logParts))
		}
	}()
}

func (s *Listener) addr() string {
	return fmt.Sprintf("%s:%d", s.host, s.port)
}

func (s *Listener) setLogFormat() {
	switch s.format {
	case RFC3164:
		s.server.SetFormat(gsyslog.RFC3164)

	case RFC6587:
		s.server.SetFormat(gsyslog.RFC6587)

	default:
		s.server.SetFormat(gsyslog.RFC5424)
	}
}

func (s *Listener) setLogProtocol() (err error) {
	if (s.proto == UDP || s.proto == TCPUDP) && err == nil {
		err = s.server.ListenUDP(s.addr())
	}

	if (s.proto == TCP || s.proto == TCPUDP) && err == nil {
		err = s.server.ListenTCP(s.addr())
	}
	return
}

func (s *Listener) Stop() {
	s.server.Kill()
}

func (s *Listener) Listen() (err error) {
	s.logpartsChannel = make(gsyslog.LogPartsChannel)
	handler := gsyslog.NewChannelHandler(s.logpartsChannel)
	server := gsyslog.NewServer()
	s.server = server
	s.setLogFormat()
	server.SetHandler(handler)

	if err = s.setLogProtocol(); err == nil {
		err = server.Boot()
	}
	return
}
