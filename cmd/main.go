package main

import (
	"my5G-RANTester/internal/control_test_engine/components/gnodeb"
	"my5G-RANTester/internal/control_test_engine/components/ue"
)

func runThreads(nUe int, nGnodeB int) {

	gnodeb.GNodeB()

	for i := 1; i <= nUe; i++ {
		go ue.Ue(i)
	}

	for {
		// do nothing
		if false {
			return
		}
	}
}
