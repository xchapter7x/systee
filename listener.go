package systee

import (
	"fmt"

	gsyslog "github.com/mcuadros/go-syslog"
	"github.com/mcuadros/go-syslog/format"
	"github.com/xchapter7x/goutil/itertools"
)

const (
	UDP = iota
	TCP
	TCPUDP
	RFC5424
	RFC3164
	RFC6587
)

type Listener struct {
	port            int
	host            string
	proto           int
	format          int
	handlerPipeline []func(LogMsg)
}

type formatter interface {
	SetFormat(format.Format)
	ListenUDP(string)
	ListenTCP(string)
}

type LogMsg map[string]interface{}

func (s *Listener) AddHandler(h func(LogMsg)) {
	handlerPipeline = append(handlerPipeline, h)
}

func (s *Listener) addr() string {
	return fmt.Sprintf("%s:%s", s.host, s.port)
}

func (s *Listener) setLogFormat(server formatter) {
	switch s.format {
	case RFC5424:
		server.SetFormat(gsyslog.RFC5424)

	case RFC3164:
		server.SetFormat(gsyslog.RFC3164)

	case RFC6587:
		server.SetFormat(gsyslog.RFC6587)
	}
}

func (s *Listener) setLogProtocol(server formatter) {
	if s.proto == UDP || s.proto == TCPUDP {
		server.ListenUDP(s.addr())
	}

	if s.proto == TCP || s.proto == TCPUDP {
		server.ListenTCP(s.addr())
	}
}

func (s *Listener) Listen() (err error) {
	channel := make(gsyslog.LogPartsChannel)
	handler := gsyslog.NewChannelHandler(channel)
	server := gsyslog.NewServer()
	s.setLogFormat(server)
	server.SetHandler(handler)
	s.setLogProtocol(server)
	server.Boot()

	go func(channel gsyslog.LogPartsChannel) {

		for logParts := range channel {
			itertools.Each(handlerPipeline, func(handler func(LogMsg)) {
				handler(logParts)
			})
		}
	}(channel)
	server.Wait()
}
