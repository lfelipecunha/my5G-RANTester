/**
*	This code is not functional.
*	It is just an example of how could we organize the UE code.
*	The base for the server-client UDS can be found in: https://gist.github.com/hakobe/6f70d69b8c5243117787fd488ae7fbf2
 */



const CM_IDLE = 0x00
const CM_REGISTERED = 0x01
const RM_DEREGISTERED = 0x00
const RM_REGISTERED = 0x01

type StateMachine struct {
	CM int
	RM int
}

func ueListen(r io.Reader) {
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf[:])
		if err != nil {
			return
		}

		// 1. Identify message
		// 2. React to message
		println("Client got:", string(buf[0:n]))
	}
}

func UserEquipment(){
	var state StateMachine
	state.CM = CM_IDLE
	state.RM = RM_DEREGISTERED

	c, err := net.Dial("unix", /tmp/my5g-gnodeb.sock)

	if err != nil {
		log.Fatal("Dial error", err)
	}
	defer c.Close()

	go ueListen(c)
	for {
		msg := "hi"
		_, err := c.Write([]byte(msg))
		if err != nil {
			log.Fatal("Write error:", err)
			break
		}
		println("Client sent:", msg)
		time.Sleep(1e9)
	}
}

func ueListen(){
	// listen
}