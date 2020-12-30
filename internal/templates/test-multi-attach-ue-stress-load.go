package templates

import (
	"my5G-RANTester/config"
	control_test_engine "my5G-RANTester/internal/control_test_engine"
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
	imsi_control := 1
	for amount <= end {
		log.Info("Launching ", amount, " UEs")
		for j := 0; j < amount; j++ {
			imsi := control_test_engine.ImsiGenerator(imsi_control)
			wg.Add(1)
			go attachUeWithTnla(imsi, cfg, int64(imsi_control), &wg)
			imsi_control++
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
	imsi_control := 1
	for amount >= end {
		for j := 0; j < amount; j++ {
			imsi := control_test_engine.ImsiGenerator(imsi_control)
			wg.Add(1)
			go attachUeWithTnla(imsi, cfg, int64(imsi_control), &wg)
			imsi_control++
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
	imsi_control := 1
	for i := 0; i <= step; i++ {
		for j := 0; j < amount; j++ {
			imsi := control_test_engine.ImsiGenerator(imsi_control)
			wg.Add(1)
			go attachUeWithTnla(imsi, cfg, int64(imsi_control), &wg)
			imsi_control++
		}
		time.Sleep(time.Duration(interval) * time.Millisecond)
	}

	// wait for multiple goroutines.
	wg.Wait()
}
