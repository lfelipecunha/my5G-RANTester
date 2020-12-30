package templates

import (
	"my5G-RANTester/config"
	control_test_engine "my5G-RANTester/internal/control_test_engine"
	"my5G-RANTester/internal/data_test_engine"
	"my5G-RANTester/internal/sctp"

	log "github.com/sirupsen/logrus"
)

// testing attach and ping for multiple queued UEs.
func TestMultiAttachUesInQueue(numberUes int) {
	var cfg = config.Data
	// UEs info
	var mcc = "208"
	var mnc = "93"

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

	/**
	** Starting message flow
	**/

	// authentication to a GNB.
	contextGnb, err := control_test_engine.RegistrationGNB(conn, cfg.GNodeB.PlmnList.GnbId, "my5GRANTester", cfg)
	if err != nil {
		log.Fatal("The test failed when GNB tried to attach! Error: ", err)
	}

	// authentication and ping to some UEs.
	for i := 1; i <= numberUes; i++ {

		// generating some IMSIs to each UE.
		imsi := control_test_engine.ImsiGenerator(i)

		wrapper := sctp.SCTPWrapper{Conn: conn}

		ue, err := control_test_engine.RegistrationUE(&wrapper, imsi, int64(i), cfg, contextGnb, mcc, mnc)
		if err != nil {
			log.Error("The test failed when UE ", imsi, " tried to attach! Error: ", err)
		}

		// data plane UE
		gtpHeader := data_test_engine.GenerateGtpHeader(int(ue.GetUeTeid()))

		err = data_test_engine.PingUE(upfConn, gtpHeader, ue.GetIp(), cfg.Ue.Ping)
		if err != nil {
			log.Error("The test failed when UE tried to use ping! Error: ", err)
		}
	}

	// end sockets.
	conn.Close()
	upfConn.Close()

}
