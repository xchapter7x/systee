package systee_test

import (
	"fmt"
	"log/syslog"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xchapter7x/systee"
)

var _ = Describe("Listener", func() {
	Describe("Integration test for a Listener who is listening", func() {
		var (
			host     string = "127.0.0.1"
			port     int    = 51222
			proto    int
			format   int
			listener *Listener
			slog     *syslog.Writer
		)

		BeforeEach(func() {
			proto = TCP
			format = RFC5424
			listener = NewListener(host, port, proto, format)
			listener.Listen()
			active := listener.IsListening()
			select {
			case <-active:
				h := fmt.Sprintf("%s:%d", host, port)
				slog, _ = syslog.Dial("tcp", h, syslog.LOG_DEBUG, "TestSyslog")
			}
		})

		AfterEach(func() {
			slog.Close()
			listener.Stop()
		})

		Context("Sending a log message over TCP", func() {
			It("should fire the handler", func() {
				var logMsg LogMsg
				var called chan bool
				var cnt int
				called = make(chan bool, 1)
				listener.AddHandler(func(lm LogMsg) {
					called <- true
					logMsg = lm
					cnt++
				})
				slog.Info("hello there")
				select {
				case <-called:
					Ω(logMsg).ShouldNot(BeNil())
					Ω(cnt).Should(Equal(1))
				}
			})
		})
	})
})
