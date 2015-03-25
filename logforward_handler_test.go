package systee_test

import (
	"math/rand"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xchapter7x/systee"
)

var _ = Describe("Logforward handler", func() {
	Describe("Logforward", func() {
		Describe("Handle", func() {
			Context("is able to connect to remote", func() {
				var (
					err      error
					host     string = "127.0.0.1"
					lh       LogHandler
					listener *Listener
					port     = 60000
					cnt      = 0
				)

				BeforeEach(func() {
					rand.Seed(time.Now().Unix())
					port++
					proto := TCP
					format := RFC5424
					listener = NewListener(host, port, proto, format)

					if e := listener.Listen(); e == nil {
						lh, err = NewLogforward(host, port, proto)
					} else {
						Ω(e).Should(BeNil())
					}

					cnt = 0
					arg := LogMsg{"hi": "Test"}
					wg := new(sync.WaitGroup)
					wg.Add(1)
					listener.AddHandler(func(lm LogMsg) {
						defer wg.Done()
						cnt++
					})
					lh.Handle(arg)
					wg.Wait()
				})

				AfterEach(func() {
					listener.Stop()
				})

				It("Should have actually written to the remote syslog", func() {
					Ω(cnt).Should(BeNumerically(">", 0))
				})
			})
		})

		Describe("Connect", func() {
			Context("can not connect to remote", func() {
				var (
					err error
					lh  LogHandler
				)

				BeforeEach(func() {
					lh, err = NewLogforward("0.1.2.3", 8000, UDP)
				})

				It("Should return a connection error", func() {
					Ω(err).ShouldNot(BeNil())
					Ω(err).Should(Equal(ErrSyslogConnectionFailure))
				})

				It("Should not return a valid log handler object", func() {
					Ω(lh).Should(BeNil())
				})
			})

			Context("is able to connect to remote", func() {
				var (
					err      error
					host     string = "127.0.0.1"
					lh       LogHandler
					listener *Listener
					port     = 50000
				)

				BeforeEach(func() {
					rand.Seed(time.Now().Unix())
					port++
					proto := TCP
					format := RFC5424
					listener = NewListener(host, port, proto, format)

					if e := listener.Listen(); e == nil {
						lh, err = NewLogforward(host, port, proto)
					} else {
						Ω(e).Should(BeNil())
					}
				})

				AfterEach(func() {
					listener.Stop()
				})

				It("Should not return a connection error", func() {
					Ω(err).Should(BeNil())
					Ω(err).ShouldNot(Equal(ErrSyslogConnectionFailure))
				})

				It("Should return a valid log handler object", func() {
					Ω(lh).ShouldNot(BeNil())
				})
			})
		})
	})
})
