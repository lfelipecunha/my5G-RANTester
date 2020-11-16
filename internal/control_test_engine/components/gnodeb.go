/**
*	This code is not functional.
*	It is just an example of how could we organize the UE code.
*	The base for the server-client UDS can be found in: https://gist.github.com/hakobe/6f70d69b8c5243117787fd488ae7fbf2
 */

func gnbListen(c net.Conn) {
	for {
		buf := make([]byte, 512)
		nr, err := c.Read(buf)
		if err != nil {
			return
		}

		data := buf[0:nr]
		println("Server got:", string(data))
		_, err = c.Write(data)
		if err != nil {
			log.Fatal("Writing client error: ", err)
		}
	}
}

func GNodeB() {
	log.Println("Starting echo server")
	ln, err := net.Listen("unix", "/tmp/go.sock")
	if err != nil {
		log.Fatal("Listen error: ", err)
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
	go func(ln net.Listener, c chan os.Signal) {
		sig := <-c
		log.Printf("Caught signal %s: shutting down.", sig)
		ln.Close()
		os.Exit(0)
	}(ln, sigc)

	for {
		fd, err := ln.Accept()
		if err != nil {
			log.Fatal("Accept error: ", err)
		}

		go echoServer(fd)
	}
}
