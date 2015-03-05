package systee_test

import (
	"fmt"
	"log/syslog"
	"sync"

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
			h := fmt.Sprintf("%s:%d", host, port)
			slog, _ = syslog.Dial("tcp", h, syslog.LOG_DEBUG, "TestSyslog")
		})

		AfterEach(func() {
			slog.Close()
			listener.Stop()
		})

		Context("Sending a log message over TCP", func() {
			var (
				logMsg LogMsg
				cnt    int = 0
				wg     *sync.WaitGroup
			)

			BeforeEach(func() {
				wg = new(sync.WaitGroup)
				wg.Add(1)
				listener.AddHandler(func(lm LogMsg) {
					defer wg.Done()
					logMsg = lm
					cnt++
				})
			})

			AfterEach(func() {
				cnt = 0
			})

			Context("single handler", func() {
				BeforeEach(func() {
					slog.Info("hello there")
					wg.Wait()
				})

				It("should fire one handler which recieves the logmsg", func() {
					Ω(logMsg).ShouldNot(BeNil())
					Ω(cnt).Should(Equal(1))
				})
			})

			//Context("multiple handlers", func() {
			//BeforeEach(func() {
			//wg.Add(1)
			//listener.AddHandler(func(lm LogMsg) {
			//defer wg.Done()
			//logMsg = lm
			//cnt++
			//})
			//slog.Info("hello there again")
			//wg.Wait()
			//})

			//It("Should fire nultiple handlers which all recieve logmsg", func() {
			//Ω(logMsg).ShouldNot(BeNil())
			//Ω(cnt).Should(BeNumerically(">", 1))
			//})
			//})
		})
	})
})
