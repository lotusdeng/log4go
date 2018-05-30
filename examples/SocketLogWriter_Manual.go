package main

import (
	"time"
)

import l4g "code.google.com/p/log4go"

func main() {
	log := l4g.NewLogger()
	log.AddFilter("network", l4g.FINEST, l4g.NewSocketLogWriter("udp", "192.168.1.255:12124"))

	// Run `nc -u -l -p 12124` or similar before you run this to see the following message
	log.Info("The time is now: %s", time.Now().Format("2006-01-02T15:04:05.123"))

	// This makes sure the output stream buffer is written
	log.Close()
}
