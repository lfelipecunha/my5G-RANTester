package ue

import (
	log "github.com/sirupsen/logrus"
	"io"
	"my5G-RANTester/config"
	"my5G-RANTester/internal/control_test_engine"
	"my5G-RANTester/internal/control_test_engine/context"
	"my5G-RANTester/internal/control_test_engine/nas_control/mm_5gs"
	"my5G-RANTester/lib/nas"
	"my5G-RANTester/lib/nas/nasMessage"
	"my5G-RANTester/lib/nas/security"
	"net"
)

var mcc = "208"
var mnc = "93"
var op = "c9e8763286b5b9ffbdf56e1297d0887b"

const sock = "/tmp/rantester.sock"

// 5GMM main states in the UE
const MM5G_DEREGISTERED = 0x00
const MM5G_REGISTERED_INITIATED = 0x01
const MM5G_REGISTERED = 0x02
const MM5G_SERVICE_REQ_INIT = 0x03
const MM5G_DEREGISTERED_INIT = 0x04

func Ue(index int) {
	var imsi = control_test_engine.ImsiGenerator(index)
	var ranUeId = int64(index)
	var conf = config.Data
	var ueContext = &context.RanUeContext{}

	ueContext.NewRanUeContext(
		imsi,
		ranUeId,
		security.AlgCiphering128NEA0,
		security.AlgIntegrity128NIA2,
		conf.Ue.Key,
		conf.Ue.Opc,
		op,
		conf.Ue.Amf,
		mcc,
		mnc,
		int32(conf.Ue.Snssai.Sd),
		conf.Ue.Snssai.Sst)

	log.Info("-------------------------------------")
	log.Info("UE INFORMATION:")
	log.Info("1-IMSI: ", ueContext.Supi)
	log.Info("2-OPc: ", conf.Ue.Opc)
	log.Info("2-Key: ", conf.Ue.Key)
	log.Info("3-Amf: ", conf.Ue.Amf)
	log.Info("-------------------------------------")

	connect(ueContext)
}

func connect(ueContext *context.RanUeContext) {
	var msg []byte
	var c net.Conn
	var ueStateMachine = MM5G_DEREGISTERED
	var receivedMsgType uint8
	var channel = make(chan uint8)

	log.Info("Connecting...")
	c, err := net.Dial("unix", sock)

	if err != nil {
		log.Fatal("Dial error", err)
	}
	defer c.Close()

	log.Info("Connected!")

	//go reader(c, channel)
	for {
		switch ueStateMachine {
		case MM5G_DEREGISTERED:
			registrationRequest := mm_5gs.GetRegistrationRequestWith5GMM(
				nasMessage.RegistrationType5GSInitialRegistration,
				ueContext.Suci,
				nil,
				nil,
				ueContext)

			msg = registrationRequest
			ueStateMachine = MM5G_REGISTERED_INITIATED

		case MM5G_REGISTERED_INITIATED:
			receivedMsgType = <-channel

			switch receivedMsgType {
			case nas.MsgTypeAuthenticationRequest:
				//Implement auth response
			case nas.MsgTypeSecurityModeCommand:
				//Implement security mode complete
			default:
				log.Warn("Unknown received message type")
			}
		case MM5G_REGISTERED:
		case MM5G_SERVICE_REQ_INIT:
		case MM5G_DEREGISTERED_INIT:
		default:
			log.Error("State not defined")
		}
		log.WithFields(log.Fields{
			"state": ueStateMachine,
		}).Info("Sending message")

		log.Info("Message size in bytes: ", len(msg))
		_, err := c.Write(msg)
		if err != nil {
			log.Fatal("Write error:", err)
		}
		log.Info("Client sent:", msg)
	}
}

func reader(r io.Reader, channel chan uint8) {
	var nasMessage nas.Message

	buf := make([]byte, 1024)
	for {
		_, err := r.Read(buf[:])
		if err != nil {
			log.Fatal("UE connection reading error")
		}
		nasMessage.PlainNasDecode(&buf)

		if nasMessage.GmmMessage != nil {
			channel <- nasMessage.GmmMessage.GmmHeader.GetMessageType()
		} else {
			log.Warn("GSM Message received but not implemented yet")
		}
	}
}

func sendMessage(conn net.Conn, msg []byte) {

	log.Info("Message size in bytes: ", len(msg))
	_, err := conn.Write(msg)
	if err != nil {
		log.Fatal("Write error:", err)
	}
	log.Info("Client sent:", msg)
}

func dial(sock string) net.Conn {
	log.Info("Connecting...")
	c, err := net.Dial("unix", sock)

	if err != nil {
		log.Fatal("Dial error", err)
	}
	defer c.Close()

	log.Info("Connected!")
	return c
}
