package control_test_engine

import (
	"fmt"
	"my5G-RANTester/lib/ngap/ngapSctp"
	"net"

	"github.com/ishidawataru/sctp"
)

func getNgapIp(amfIP string, amfPort int) (amfAddr *sctp.SCTPAddr, err error) {
	ips := []net.IPAddr{}
	// se der um erro != nill entra no if.
	if ip, err1 := net.ResolveIPAddr("ip", amfIP); err1 != nil {
		err = fmt.Errorf("Error resolving address '%s': %v", amfIP, err1)
		return
	} else {
		ips = append(ips, *ip)
	}
	amfAddr = &sctp.SCTPAddr{
		IPAddrs: ips,
		Port:    amfPort,
	}
	return
}

func ConnectToAmf(amfIP string, amfPort int) (*sctp.SCTPConn, error) {
	amfAddr, err := getNgapIp(amfIP, amfPort)
	if err != nil {
		return nil, err
	}
	conn, err := sctp.DialSCTP("sctp", nil, amfAddr)
	if err != nil {
		return nil, err
	}
	info, _ := conn.GetDefaultSentParam()
	info.PPID = ngapSctp.NGAP_PPID
	err = conn.SetDefaultSentParam(info)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
