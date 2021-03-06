package mm_5gs

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"my5G-RANTester/internal/control_test_engine/context"
	"my5G-RANTester/internal/control_test_engine/nas_control"
	"my5G-RANTester/lib/nas"
	"my5G-RANTester/lib/nas/nasMessage"
	"my5G-RANTester/lib/nas/nasType"
	"my5G-RANTester/lib/ngap/ngapType"
)

func getAuthenticationResponse(authenticationResponseParam []uint8, eapMsg string) (nasPdu []byte) {

	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeAuthenticationResponse)

	authenticationResponse := nasMessage.NewAuthenticationResponse(0)
	authenticationResponse.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	authenticationResponse.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	authenticationResponse.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(0)
	authenticationResponse.AuthenticationResponseMessageIdentity.SetMessageType(nas.MsgTypeAuthenticationResponse)

	if len(authenticationResponseParam) > 0 {
		authenticationResponse.AuthenticationResponseParameter = nasType.NewAuthenticationResponseParameter(nasMessage.AuthenticationResponseAuthenticationResponseParameterType)
		authenticationResponse.AuthenticationResponseParameter.SetLen(uint8(len(authenticationResponseParam)))
		copy(authenticationResponse.AuthenticationResponseParameter.Octet[:], authenticationResponseParam[0:16])
	} else if eapMsg != "" {
		rawEapMsg, _ := base64.StdEncoding.DecodeString(eapMsg)
		authenticationResponse.EAPMessage = nasType.NewEAPMessage(nasMessage.AuthenticationResponseEAPMessageType)
		authenticationResponse.EAPMessage.SetLen(uint16(len(rawEapMsg)))
		authenticationResponse.EAPMessage.SetEAPMessage(rawEapMsg)
	}

	m.GmmMessage.AuthenticationResponse = authenticationResponse

	data := new(bytes.Buffer)
	err := m.GmmMessageEncode(data)
	if err != nil {
		fmt.Println(err.Error())
	}

	nasPdu = data.Bytes()
	return
}

func AuthenticationResponse(ue *context.RanUeContext, ngapMsg *ngapType.NGAPPDU) ([]byte, error) {

	// Calculate for RES*
	nasPdu := nas_control.GetNasPduFromDownlink(ngapMsg.InitiatingMessage.Value.DownlinkNASTransport)
	if nasPdu == nil {
		return nil, fmt.Errorf("Error in get NAS PDU from Authentication request")
	}

	// get rand and autn from the message.
	rand := nasPdu.AuthenticationRequest.GetRANDValue()
	autn := nasPdu.AuthenticationRequest.GetAUTN()

	resStat := ue.DeriveRESstarAndSetKey(ue.AuthenticationSubs, rand[:], ue.Snn, autn[:])

	// send NAS Authentication Response.
	pdu := getAuthenticationResponse(resStat, "")
	return pdu, nil
}
