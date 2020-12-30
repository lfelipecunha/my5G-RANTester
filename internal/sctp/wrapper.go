package sctp

import (
	"time"

	"github.com/ishidawataru/sctp"
	log "github.com/sirupsen/logrus"
)

type SCTPWrapper struct {
	Conn *sctp.SCTPConn
}

func (sw *SCTPWrapper) Write(b []byte) (int, error) {
	t1 := time.Now()

	n, err := sw.Conn.Write(b)

	t2 := time.Now()
	total := t2.UnixNano() - t1.UnixNano()

	log.Infof("Write in: %d ns", total)

	return n, err
}

func (sw *SCTPWrapper) Read(dst []byte) (int, error) {
	t1 := time.Now()

	n, err := sw.Conn.Read(dst)

	t2 := time.Now()
	total := t2.UnixNano() - t1.UnixNano()

	log.Infof("Read in: %d ns", total)

	return n, err
}
