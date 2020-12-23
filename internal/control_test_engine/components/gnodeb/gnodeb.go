package gnodeb

import (
	log "github.com/sirupsen/logrus"
	"my5G-RANTester/config"
	"my5G-RANTester/internal/control_test_engine"
	"my5G-RANTester/internal/control_test_engine/context"
	"my5G-RANTester/internal/data_test_engine"
	"my5G-RANTester/lib/nas"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const sock = "/tmp/rantester.sock"

func GNodeB() {
	var cfg = config.Data

	log.Info("[CORE]", cfg.AMF.Name, " Core in Testing")
	log.Info("Conecting to AMF...")
	conn, err := control_test_engine.ConnectToAmf(cfg.AMF.Ip, cfg.GNodeB.ControlIF.Ip, cfg.AMF.Port, cfg.GNodeB.ControlIF.Port)
	if err != nil {
		log.Fatal("The test failed when sctp socket tried to connect to AMF! Error: ", err)
	}
	log.Info("OK")

	log.Info("Conecting to UPF...")
	upfConn, err := data_test_engine.ConnectToUpf(cfg.GNodeB.DataIF.Ip, cfg.UPF.Ip, cfg.GNodeB.DataIF.Port, cfg.UPF.Port)
	if err != nil {
		log.Fatal("The test failed when udp socket tried to connect to UPF! Error: ", err)
	}
	log.Info("OK")

	contextGnb, err := control_test_engine.RegistrationGNB(conn, cfg.GNodeB.PlmnList.GnbId, "my5GRANTester", cfg)
	if err != nil {
		log.Fatal("The test failed when GNB tried to attach! Error: ", err)
	}

	go startServer(upfConn, contextGnb)
}

func startServer(upfConn *net.UDPConn, ranGnbContext *context.RanGnbContext) {
	log.Info("Starting gNodeB server")
	ln, err := net.Listen("unix", sock)
	if err != nil {
		log.Fatal("Listen error: ", err)
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
	go func(ln net.Listener, c chan os.Signal) {
		sig := <-c
		log.Info("Caught signal %s: shutting down.", sig)
		ln.Close()
		os.Exit(0)
	}(ln, sigc)

	for {
		log.Info("Server waiting to accept connection")
		fd, err := ln.Accept()
		if err != nil {
			log.Fatal("Accept error: ", err)
		}

		go ueConnection(fd)
	}
}

func ueConnection(c net.Conn) {
	var nasMessage nas.Message

	buf := make([]byte, 1024)
	for {
		_, err := c.Read(buf)
		if err != nil {
			log.Fatal("gNodeB connection reading error")
		}

		nasMessage.PlainNasDecode(&buf)

		if nasMessage.GmmMessage != nil {
			log.Info("Message type: ", nasMessage.GmmMessage.GmmHeader.GetMessageType())
		} else {
			log.Warn("GSM Message received but not implemented yet")
		}
	}
}
