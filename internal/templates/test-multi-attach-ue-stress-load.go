package templates

import (
	"fmt"
	"my5G-RANTester/config"
	control_test_engine "my5G-RANTester/internal/control_test_engine"
	"my5G-RANTester/internal/sctp"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// testing attach and ping for multiple concurrent UEs using TNLAs.
func TestMultiAttachUesLoadStress(start int, step int, end int, interval int) {

	cfg, err := config.GetConfig()
	if err != nil {
		//return nil
		log.Fatal("Error in get configuration")
	}

	// authentication and ping to some concurrent UEs.
	log.Info("Testing Stress Load: Start[", start, "] Step[", step, "] End[", end, "] Interval[", interval, " seconds]")
	log.Info("[CORE]", cfg.AMF.Name, " Core in Testing")

	if start > end {
		testDecrease(start, end, step, interval, cfg)
		return
	} else if start < end {
		testIncrease(start, end, step, interval, cfg)
		return
	}

	testConstant(start, step, interval, cfg)

}

func testIncrease(start int, end int, step int, interval int, cfg config.Config) {
	log.Info("Starting Increase Test")

	// Launch several goroutines and increment the WaitGroup counter for each.
	var wg sync.WaitGroup
	amount := start
	imsiControl := 1
	for amount <= end {
		log.Info("Launching ", amount, " UEs")
		for j := 0; j < amount; j++ {
			imsi := control_test_engine.ImsiGenerator(imsiControl)
			wg.Add(1)
			go attachUEMultipleGNBs(imsi, cfg, int64(imsiControl), &wg)
			imsiControl++
			time.Sleep(10 * time.Millisecond)
		}
		amount += step
		time.Sleep(time.Duration(interval) * time.Millisecond)
	}

	// wait for multiple goroutines.
	wg.Wait()
}

func testDecrease(start int, end int, step int, interval int, cfg config.Config) {
	log.Info("Starting Decrease Test")
	// Launch several goroutines and increment the WaitGroup counter for each.
	var wg sync.WaitGroup
	amount := start
	imsiControl := 1
	for amount >= end {
		for j := 0; j < amount; j++ {
			imsi := control_test_engine.ImsiGenerator(imsiControl)
			wg.Add(1)
			go attachUEMultipleGNBs(imsi, cfg, int64(imsiControl), &wg)
			imsiControl++
		}
		amount -= step
		time.Sleep(time.Duration(interval) * time.Millisecond)
	}

	// wait for multiple goroutines.
	wg.Wait()
}

func testConstant(amount int, step int, interval int, cfg config.Config) {
	log.Info("Starting Constant Test")
	// Launch several goroutines and increment the WaitGroup counter for each.
	var wg sync.WaitGroup
	imsiControl := 1
	for i := 0; i <= step; i++ {
		for j := 0; j < amount; j++ {
			imsi := control_test_engine.ImsiGenerator(imsiControl)
			wg.Add(1)
			go attachUEMultipleGNBs(imsi, cfg, int64(imsiControl), &wg)
			imsiControl++
		}
		time.Sleep(time.Duration(interval) * time.Millisecond)
	}

	// wait for multiple goroutines.
	wg.Wait()
}

// testing attach and ping for a UE with TNLA.
func attachUEMultipleGNBs(imsi string, conf config.Config, ranUeId int64, wg *sync.WaitGroup) {

	defer wg.Done()

	// make N2(RAN connect to AMF)
	log.Info("Conecting to AMF...")
	conn, err := control_test_engine.ConnectToAmf(conf.AMF.Ip, conf.AMF.Port)
	if err != nil {
		log.Errorf("The test failed when sctp socket tried to connect to AMF! Error:", err)
		return
	}
	log.Info("OK")

	gnbID := int(ranUeId / 255)

	ranUeId = ranUeId % 255

	// multiple names for GNBs.
	gnbName := fmt.Sprint("my5gRanTester", gnbID)

	// generate GNB id.
	var gnbIDString string
	if gnbID < 16 {
		gnbIDString = "00000" + fmt.Sprintf("%x", gnbID)
	} else if gnbID < 256 {
		gnbIDString = "0000" + fmt.Sprintf("%x", gnbID)
	} else {
		gnbIDString = "000" + fmt.Sprintf("%x", gnbID)
	}

	// authentication to a GNB.
	gnbContext, err := control_test_engine.RegistrationGNB(conn, gnbIDString, gnbName, conf)
	if err != nil {
		log.Errorf("The test failed when GNB tried to attach! Error:", err)
		conn.Close()
		return
	}

	wrapper := sctp.SCTPWrapper{Conn: conn}

	ue, err := control_test_engine.RegistrationUE(&wrapper, imsi, ranUeId, conf, gnbContext, "208", "93")
	if err != nil {
		log.Error("The test failed when UE", ranUeId, " Suci: ", ue.Suci, "tried to attach! Error:", err)
	}

	// end sctp socket.
	conn.Close()

}
