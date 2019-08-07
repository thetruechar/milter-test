package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"os"

	"github.com/phalaaxx/milter"
)

type TestMilter struct {
	milter.Milter
	from string
	tos  []string
}

func (TestMilter) Connect(host string, family string, port uint16, addr net.IP, m *milter.Modifier) (milter.Response, error) {
	log.Println(fmt.Sprintf("Connect---%s,%s,%d,%v", host, family, port, addr))
	return milter.RespContinue, nil
}

func (TestMilter) Helo(name string, m *milter.Modifier) (milter.Response, error) {
	log.Println(fmt.Sprintf("Helo---name:%s", name))
	return milter.RespContinue, nil
}

func (TestMilter) MailFrom(from string, m *milter.Modifier) (milter.Response, error) {
	log.Println(fmt.Sprintf("MailFrom---from:%s", from))
	return milter.RespContinue, nil
}

func (TestMilter) RcptTo(rcptTo string, m *milter.Modifier) (milter.Response, error) {
	log.Println(fmt.Sprintf("RcptTo---rcptTo:%s", rcptTo))
	return milter.RespContinue, nil
}

func (TestMilter) Header(name string, value string, m *milter.Modifier) (milter.Response, error) {
	log.Println(fmt.Sprintf("Header---name:%s,value:%s", name, value))
	return milter.RespContinue, nil
}

func (TestMilter) Headers(h textproto.MIMEHeader, m *milter.Modifier) (milter.Response, error) {
	log.Println(fmt.Sprintf("Headers---h:%v", h, ))
	return milter.RespContinue, nil
}

func (TestMilter) BodyChunk(chunk []byte, m *milter.Modifier) (milter.Response, error) {
	log.Println(fmt.Sprintf("BodyChunk---"))
	return milter.RespContinue, nil
}

func (TestMilter) Body(m *milter.Modifier) (milter.Response, error) {
	log.Println(fmt.Sprintf("Body---"))
	return milter.RespContinue, nil
}

func RunServer(socket net.Listener) {
	log.Println("milter RunServer...")
	if err := milter.RunServer(socket,
		func() (milter.Milter, milter.OptAction, milter.OptProtocol) {
			return &TestMilter{},
				0,
				0
		}); err != nil {
		log.Fatal(err)
	}
}

func main() {
	var protocol, address string
	flag.StringVar(&protocol,
		"proto",
		"tcp",
		"Protocol family (unix or tcp)")
	flag.StringVar(&address,
		"addr",
		"127.0.0.1:9001",
		"Bind to address or unix domain socket")
	if protocol == "unix" {
		// ignore os.Remove errors
		_ = os.Remove(address)
	}
	socket, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatal(err)
	}
	defer socket.Close()

	if protocol == "unix" {
		// set mode 0660 for unix domain sockets
		if err := os.Chmod(address, 0660); err != nil {
			log.Fatal(err)
		}
		// remove socket on exit
		defer os.Remove(address)
	}
	go RunServer(socket)
	select {}
}
