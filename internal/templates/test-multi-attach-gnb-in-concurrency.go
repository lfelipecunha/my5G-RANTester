package templates

import (
	"fmt"
	"my5G-RANTester/config"
	"my5G-RANTester/internal/control_test_engine"
	"sync"

	log "github.com/sirupsen/logrus"
)

func TestMultiAttachGnbInConcurrency(numberGnbs int) {

	var wg sync.WaitGroup

	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal("Error in get configuration")
	}
	// fmt.Printf("[CORE]%s Core in Testing\n", cfg.AMF.Name)

	log.Info("Testing attached with ", numberGnbs, " gnbs")
	log.Info("[CORE]", cfg.AMF.Name, " Core in Testing")

	// multiple concurrent GNBs authentication using goroutines.
	for i := 1; i <= numberGnbs; i++ {

		wg.Add(1)
		go func(wg *sync.WaitGroup, i int) {

			defer wg.Done()

			// make N2(RAN connect to AMF)
			conn, err := control_test_engine.ConnectToAmf(cfg.AMF.Ip, cfg.AMF.Port)
			if err != nil {
				log.Fatal("The test failed when sctp socket tried to connect to AMF! Error:", err)
			}

			// multiple names for GNBs.
			nameGNB := fmt.Sprint("my5gRanTester", i)
			// fmt.Println(nameGNB)

			// generate GNB id.
			var aux string
			if i < 16 {
				aux = "00000" + fmt.Sprintf("%x", i)
			} else if i < 256 {
				aux = "0000" + fmt.Sprintf("%x", i)
			} else {
				aux = "000" + fmt.Sprintf("%x", i)
			}

			// authentication to a GNB.
			contextgnb, err := control_test_engine.RegistrationGNB(conn, aux, nameGNB, cfg)
			if err != nil || contextgnb == nil {
				log.Error("The test failed when GNB", aux, "tried to attach! Error:", err)
			}

			//fmt.Println(contextgnb)

			// close sctp socket.
			conn.Close()
		}(&wg, i)
	}

	// wait threads.
	wg.Wait()

	// return nil
}
