package pdu_session_management

import (
	"encoding/hex"
	"fmt"
	"my5G-RANTester/lib/aper"
	"my5G-RANTester/lib/ngap"
	"my5G-RANTester/lib/ngap/ngapConvert"
	"my5G-RANTester/lib/ngap/ngapType"

	"my5G-RANTester/internal/sctp"
)

func PDUSessionResourceSetupResponse(connN2 *sctp.SCTPWrapper, amfUeNgapID int64, ranUeNgapID int64, supi string, ranIpAddr string) error {
	sendMsg, err := getPDUSessionResourceSetupResponse(amfUeNgapID, ranUeNgapID, ranIpAddr)
	if err != nil {
		return fmt.Errorf("Error getting %s ue NGAP-PDU Session Resource Setup Response", supi)
	}
	_, err = connN2.Write(sendMsg)
	if err != nil {
		return fmt.Errorf("Error sending %s ue NGAP-PDU Session Resource Setup Response", supi)
	}

	return nil
}

func getPDUSessionResourceSetupResponse(amfUeNgapID int64, ranUeNgapID int64, ipv4 string) ([]byte, error) {
	message := buildPDUSessionResourceSetupResponseForRegistrationTest(amfUeNgapID, ranUeNgapID, ipv4)
	return ngap.Encoder(message)
}

func buildPDUSessionResourceSetupResponseForRegistrationTest(amfUeNgapID, ranUeNgapID int64, ipv4 string) (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentSuccessfulOutcome
	pdu.SuccessfulOutcome = new(ngapType.SuccessfulOutcome)

	successfulOutcome := pdu.SuccessfulOutcome
	successfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodePDUSessionResourceSetup
	successfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	successfulOutcome.Value.Present = ngapType.SuccessfulOutcomePresentPDUSessionResourceSetupResponse
	successfulOutcome.Value.PDUSessionResourceSetupResponse = new(ngapType.PDUSessionResourceSetupResponse)

	pDUSessionResourceSetupResponse := successfulOutcome.Value.PDUSessionResourceSetupResponse
	pDUSessionResourceSetupResponseIEs := &pDUSessionResourceSetupResponse.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.PDUSessionResourceSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PDUSessionResourceSetupResponseIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = amfUeNgapID

	pDUSessionResourceSetupResponseIEs.List = append(pDUSessionResourceSetupResponseIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.PDUSessionResourceSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PDUSessionResourceSetupResponseIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ranUeNgapID

	pDUSessionResourceSetupResponseIEs.List = append(pDUSessionResourceSetupResponseIEs.List, ie)

	// PDU Session Resource Setup Response List
	ie = ngapType.PDUSessionResourceSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceSetupListSURes
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PDUSessionResourceSetupResponseIEsPresentPDUSessionResourceSetupListSURes
	ie.Value.PDUSessionResourceSetupListSURes = new(ngapType.PDUSessionResourceSetupListSURes)

	pDUSessionResourceSetupListSURes := ie.Value.PDUSessionResourceSetupListSURes

	// PDU Session Resource Setup Response Item in PDU Session Resource Setup Response List
	pDUSessionResourceSetupItemSURes := ngapType.PDUSessionResourceSetupItemSURes{}

	// PDU Session ID : This is an unique identifier generated by UE. Can’t be same as any existing PDU session.
	pDUSessionResourceSetupItemSURes.PDUSessionID.Value = ranUeNgapID

	pDUSessionResourceSetupItemSURes.PDUSessionResourceSetupResponseTransfer = GetPDUSessionResourceSetupResponseTransfer(ipv4, amfUeNgapID)

	pDUSessionResourceSetupListSURes.List = append(pDUSessionResourceSetupListSURes.List, pDUSessionResourceSetupItemSURes)

	pDUSessionResourceSetupResponseIEs.List = append(pDUSessionResourceSetupResponseIEs.List, ie)

	// PDU Sessuin Resource Failed to Setup List
	// ie = ngapType.PDUSessionResourceSetupResponseIEs{}
	// ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceFailedToSetupListSURes
	// ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	// ie.Value.Present = ngapType.PDUSessionResourceSetupResponseIEsPresentPDUSessionResourceFailedToSetupListSURes
	// ie.Value.PDUSessionResourceFailedToSetupListSURes = new(ngapType.PDUSessionResourceFailedToSetupListSURes)

	// pDUSessionResourceFailedToSetupListSURes := ie.Value.PDUSessionResourceFailedToSetupListSURes

	// // PDU Session Resource Failed to Setup Item in PDU Sessuin Resource Failed to Setup List
	// pDUSessionResourceFailedToSetupItemSURes := ngapType.PDUSessionResourceFailedToSetupItemSURes{}
	// pDUSessionResourceFailedToSetupItemSURes.PDUSessionID.Value = 10
	// pDUSessionResourceFailedToSetupItemSURes.PDUSessionResourceSetupUnsuccessfulTransfer = GetPDUSessionResourceSetupUnsucessfulTransfer()

	// pDUSessionResourceFailedToSetupListSURes.List = append(pDUSessionResourceFailedToSetupListSURes.List, pDUSessionResourceFailedToSetupItemSURes)

	// pDUSessionResourceSetupResponseIEs.List = append(pDUSessionResourceSetupResponseIEs.List, ie)
	// Criticality Diagnostics (optional)
	return
}

func GetPDUSessionResourceSetupResponseTransfer(ipv4 string, amfId int64) []byte {
	data := buildPDUSessionResourceSetupResponseTransfer(ipv4, amfId)
	encodeData, _ := aper.MarshalWithParams(data, "valueExt")
	return encodeData
}

func buildPDUSessionResourceSetupResponseTransfer(ipv4 string, amfId int64) (data ngapType.PDUSessionResourceSetupResponseTransfer) {

	// QoS Flow per TNL Information
	qosFlowPerTNLInformation := &data.QosFlowPerTNLInformation
	qosFlowPerTNLInformation.UPTransportLayerInformation.Present = ngapType.UPTransportLayerInformationPresentGTPTunnel

	// UP Transport Layer Information in QoS Flow per TNL Information
	upTransportLayerInformation := &qosFlowPerTNLInformation.UPTransportLayerInformation
	upTransportLayerInformation.Present = ngapType.UPTransportLayerInformationPresentGTPTunnel
	upTransportLayerInformation.GTPTunnel = new(ngapType.GTPTunnel)

	// generates some GTP-TEIDs for UPF-RAN tunnels(downlink)
	var aux string
	if amfId < 16 {
		aux = "0000000" + fmt.Sprintf("%x", amfId)
	} else if amfId < 256 {
		aux = "000000" + fmt.Sprintf("%x", amfId)
	} else {
		aux = "00000" + fmt.Sprintf("%x", amfId)
	}
	resu, err := hex.DecodeString(aux)
	if err != nil {
		fmt.Println("error in GTPTEID for endpoint UPF-RAN")
		fmt.Println(err)
	}
	upTransportLayerInformation.GTPTunnel.GTPTEID.Value = aper.OctetString(resu)
	upTransportLayerInformation.GTPTunnel.TransportLayerAddress = ngapConvert.IPAddressToNgap(ipv4, "")

	// Associated QoS Flow List in QoS Flow per TNL Information
	associatedQosFlowList := &qosFlowPerTNLInformation.AssociatedQosFlowList

	associatedQosFlowItem := ngapType.AssociatedQosFlowItem{}
	associatedQosFlowItem.QosFlowIdentifier.Value = 1
	associatedQosFlowList.List = append(associatedQosFlowList.List, associatedQosFlowItem)

	return
}
