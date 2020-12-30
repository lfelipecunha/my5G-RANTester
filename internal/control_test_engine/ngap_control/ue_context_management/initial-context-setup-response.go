package ue_context_management

import (
	"fmt"
	"my5G-RANTester/internal/sctp"
	"my5G-RANTester/lib/ngap"
	"my5G-RANTester/lib/ngap/ngapType"
)

func InitialContextSetupResponse(connN2 *sctp.SCTPWrapper, amfUeNgapID int64, ranUeNgapID int64, supi string) error {

	sendMsg, err := getInitialContextSetupResponse(amfUeNgapID, ranUeNgapID)
	if err != nil {
		return fmt.Errorf("Error getting %s ue ngap Initial Context Setup Response Msg", supi)
	}
	_, err = connN2.Write(sendMsg)
	if err != nil {
		return fmt.Errorf("Error sending %s ue Initial Context Setup Response Msg", supi)
	}

	return nil
}

func getInitialContextSetupResponse(amfUeNgapID int64, ranUeNgapID int64) ([]byte, error) {
	message := BuildInitialContextSetupResponseForRegistraionTest(amfUeNgapID, ranUeNgapID)

	return ngap.Encoder(message)
}

func BuildInitialContextSetupResponseForRegistraionTest(amfUeNgapID, ranUeNgapID int64) (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentSuccessfulOutcome
	pdu.SuccessfulOutcome = new(ngapType.SuccessfulOutcome)

	successfulOutcome := pdu.SuccessfulOutcome
	successfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodeInitialContextSetup
	successfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	successfulOutcome.Value.Present = ngapType.SuccessfulOutcomePresentInitialContextSetupResponse
	successfulOutcome.Value.InitialContextSetupResponse = new(ngapType.InitialContextSetupResponse)

	initialContextSetupResponse := successfulOutcome.Value.InitialContextSetupResponse
	initialContextSetupResponseIEs := &initialContextSetupResponse.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.InitialContextSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.InitialContextSetupResponseIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = amfUeNgapID

	initialContextSetupResponseIEs.List = append(initialContextSetupResponseIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.InitialContextSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.InitialContextSetupResponseIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ranUeNgapID

	initialContextSetupResponseIEs.List = append(initialContextSetupResponseIEs.List, ie)

	return
}
