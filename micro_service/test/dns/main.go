package main

import (
	"net"
	"time"

	"github.com/miekg/dns"
	"github.com/sirupsen/logrus"
)

func main() {
	config, _ := dns.ClientConfigFromFile("/etc/resolv.conf")
	c := new(dns.Client)
	logrus.Infoln(config.Port)
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn("userService.service.consul"), dns.TypeNone)
	m.RecursionDesired = true
	r, _, err := c.Exchange(m, net.JoinHostPort("172.16.7.120", "8600"))
	if err != nil {
		logrus.Infoln(err)
	}
	logrus.Infoln(r.Answer)
	conn, err := net.DialTimeout("tcp", "userService.service.consul:5000", 5*time.Second)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Infoln("connect success")
	conn.Close()
}
