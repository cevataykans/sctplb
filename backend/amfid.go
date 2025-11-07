package backend

import (
	"fmt"
	"github.com/omec-project/ngap/logger"
	"github.com/omec-project/ngap/ngapType"
	"github.com/omec-project/sctplb/context"
)

func findBackendWithNGAPID(ctx *context.SctplbContext, ngapId *ngapType.AMFUENGAPID) (Backend, error) {
	if ngapId == nil {
		return nil, fmt.Errorf("ngapId is nil")
	}

	id, err := drsmClient.FindOwnerInt32ID(int32(ngapId.Value))
	if err != nil {
		return nil, err
	}
	if id == nil {
		return nil, fmt.Errorf("ngapId not found by DRSM")
	}
	logger.NgapLog.Infoln("Found backend with id:", id)

	for _, instance := range ctx.Backends {
		b1 := instance.(*GrpcServer)
		// AMF sets RedirectID as PodIp
		if b1.address != id.PodIp {
			continue
		}

		// We have the correct AMF instance
		if !b1.state {
			return nil, fmt.Errorf("backend found but not ready")
		}
		return instance, nil
	}
	return nil, fmt.Errorf("backend not found")
}

// Taken and adapted from ngap/handler.go
func ExtractAMFUENGAPID(message *ngapType.NGAPPDU) *ngapType.AMFUENGAPID {
	var aMFUENGAPID *ngapType.AMFUENGAPID

	if message == nil {
		logger.NgapLog.Errorln("NGAP Message is nil")
		return nil
	}

	switch message.Present {
	case ngapType.NGAPPDUPresentInitiatingMessage:
		initiatingMessage := message.InitiatingMessage
		if initiatingMessage == nil {
			logger.NgapLog.Errorln("initiatingMessage is nil")
			return nil
		}
		switch initiatingMessage.ProcedureCode.Value {
		case ngapType.ProcedureCodeNGSetup:
		case ngapType.ProcedureCodeInitialUEMessage:
			// no references to the aMFUENGAPID which is required by the DRSM module.

		case ngapType.ProcedureCodeUplinkNASTransport:
			ngapMsg := initiatingMessage.Value.UplinkNASTransport
			if ngapMsg == nil {
				logger.NgapLog.Errorln("UplinkNasTransport is nil")
				return nil
			}
			for i := 0; i < len(ngapMsg.ProtocolIEs.List); i++ {
				ie := ngapMsg.ProtocolIEs.List[i]
				if ie.Id.Value == ngapType.ProtocolIEIDAMFUENGAPID {
					aMFUENGAPID = ie.Value.AMFUENGAPID
				}
			}

		case ngapType.ProcedureCodeHandoverCancel:
			ngapMsg := initiatingMessage.Value.HandoverCancel
			for i := 0; i < len(ngapMsg.ProtocolIEs.List); i++ {
				ie := ngapMsg.ProtocolIEs.List[i]
				if ie.Id.Value == ngapType.ProtocolIEIDAMFUENGAPID {
					aMFUENGAPID = ie.Value.AMFUENGAPID
				}
			}

		case ngapType.ProcedureCodeUEContextReleaseRequest:
			ngapMsg := initiatingMessage.Value.UEContextReleaseRequest
			for i := 0; i < len(ngapMsg.ProtocolIEs.List); i++ {
				ie := ngapMsg.ProtocolIEs.List[i]
				if ie.Id.Value == ngapType.ProtocolIEIDAMFUENGAPID {
					aMFUENGAPID = ie.Value.AMFUENGAPID
				}
			}

		case ngapType.ProcedureCodeNASNonDeliveryIndication:
			ngapMsg := initiatingMessage.Value.NASNonDeliveryIndication
			for i := 0; i < len(ngapMsg.ProtocolIEs.List); i++ {
				ie := ngapMsg.ProtocolIEs.List[i]
				if ie.Id.Value == ngapType.ProtocolIEIDAMFUENGAPID {
					aMFUENGAPID = ie.Value.AMFUENGAPID
				}
			}

		case ngapType.ProcedureCodeLocationReportingFailureIndication:
		case ngapType.ProcedureCodeErrorIndication:
		case ngapType.ProcedureCodeUERadioCapabilityInfoIndication:
			ngapMsg := initiatingMessage.Value.UERadioCapabilityInfoIndication
			for i := 0; i < len(ngapMsg.ProtocolIEs.List); i++ {
				ie := ngapMsg.ProtocolIEs.List[i]
				if ie.Id.Value == ngapType.ProtocolIEIDAMFUENGAPID {
					aMFUENGAPID = ie.Value.AMFUENGAPID
				}
			}

		case ngapType.ProcedureCodeHandoverNotification:
			ngapMsg := initiatingMessage.Value.HandoverNotify
			for i := 0; i < len(ngapMsg.ProtocolIEs.List); i++ {
				ie := ngapMsg.ProtocolIEs.List[i]
				if ie.Id.Value == ngapType.ProtocolIEIDAMFUENGAPID {
					aMFUENGAPID = ie.Value.AMFUENGAPID
				}
			}

		case ngapType.ProcedureCodeHandoverPreparation:
			ngapMsg := initiatingMessage.Value.HandoverRequired
			for i := 0; i < len(ngapMsg.ProtocolIEs.List); i++ {
				ie := ngapMsg.ProtocolIEs.List[i]
				if ie.Id.Value == ngapType.ProtocolIEIDAMFUENGAPID {
					aMFUENGAPID = ie.Value.AMFUENGAPID
				}
			}

		case ngapType.ProcedureCodeRANConfigurationUpdate:
		case ngapType.ProcedureCodeRRCInactiveTransitionReport:
		case ngapType.ProcedureCodePDUSessionResourceNotify:
			ngapMsg := initiatingMessage.Value.PDUSessionResourceNotify
			for i := 0; i < len(ngapMsg.ProtocolIEs.List); i++ {
				ie := ngapMsg.ProtocolIEs.List[i]
				if ie.Id.Value == ngapType.ProtocolIEIDAMFUENGAPID {
					aMFUENGAPID = ie.Value.AMFUENGAPID
				}
			}

		case ngapType.ProcedureCodePathSwitchRequest:
			ngapMsg := initiatingMessage.Value.PathSwitchRequest
			for i := 0; i < len(ngapMsg.ProtocolIEs.List); i++ {
				ie := ngapMsg.ProtocolIEs.List[i]
				if ie.Id.Value == ngapType.ProtocolIEIDSourceAMFUENGAPID {
					aMFUENGAPID = ie.Value.SourceAMFUENGAPID
				}
			}

		case ngapType.ProcedureCodeLocationReport:
		case ngapType.ProcedureCodeUplinkUEAssociatedNRPPaTransport:
		case ngapType.ProcedureCodeUplinkRANConfigurationTransfer:
		case ngapType.ProcedureCodePDUSessionResourceModifyIndication:
			ngapMsg := initiatingMessage.Value.PDUSessionResourceModifyIndication
			for i := 0; i < len(ngapMsg.ProtocolIEs.List); i++ {
				ie := ngapMsg.ProtocolIEs.List[i]
				if ie.Id.Value == ngapType.ProtocolIEIDAMFUENGAPID {
					aMFUENGAPID = ie.Value.AMFUENGAPID
				}
			}

		case ngapType.ProcedureCodeCellTrafficTrace:
		case ngapType.ProcedureCodeUplinkRANStatusTransfer:
		case ngapType.ProcedureCodeUplinkNonUEAssociatedNRPPaTransport:
		}

	case ngapType.NGAPPDUPresentSuccessfulOutcome:
		successfulOutcome := message.SuccessfulOutcome
		if successfulOutcome == nil {
			logger.NgapLog.Errorln("successfulOutcome is nil")
			return nil
		}

		switch successfulOutcome.ProcedureCode.Value {
		case ngapType.ProcedureCodeNGReset:
		case ngapType.ProcedureCodeUEContextRelease:
			ngapMsg := successfulOutcome.Value.UEContextReleaseComplete
			for i := 0; i < len(ngapMsg.ProtocolIEs.List); i++ {
				ie := ngapMsg.ProtocolIEs.List[i]
				if ie.Id.Value == ngapType.ProtocolIEIDAMFUENGAPID {
					aMFUENGAPID = ie.Value.AMFUENGAPID
				}
			}

		case ngapType.ProcedureCodePDUSessionResourceRelease:
			ngapMsg := successfulOutcome.Value.PDUSessionResourceReleaseResponse
			for i := 0; i < len(ngapMsg.ProtocolIEs.List); i++ {
				ie := ngapMsg.ProtocolIEs.List[i]
				if ie.Id.Value == ngapType.ProtocolIEIDAMFUENGAPID {
					aMFUENGAPID = ie.Value.AMFUENGAPID
				}
			}

		case ngapType.ProcedureCodeUERadioCapabilityCheck:
		case ngapType.ProcedureCodeAMFConfigurationUpdate:
		case ngapType.ProcedureCodeInitialContextSetup:
			ngapMsg := successfulOutcome.Value.InitialContextSetupResponse
			for i := 0; i < len(ngapMsg.ProtocolIEs.List); i++ {
				ie := ngapMsg.ProtocolIEs.List[i]
				if ie.Id.Value == ngapType.ProtocolIEIDAMFUENGAPID {
					aMFUENGAPID = ie.Value.AMFUENGAPID
				}
			}

		case ngapType.ProcedureCodeUEContextModification:
			ngapMsg := successfulOutcome.Value.UEContextModificationResponse
			for i := 0; i < len(ngapMsg.ProtocolIEs.List); i++ {
				ie := ngapMsg.ProtocolIEs.List[i]
				if ie.Id.Value == ngapType.ProtocolIEIDAMFUENGAPID {
					aMFUENGAPID = ie.Value.AMFUENGAPID
				}
			}

		case ngapType.ProcedureCodePDUSessionResourceSetup:
			ngapMsg := successfulOutcome.Value.PDUSessionResourceSetupResponse
			for i := 0; i < len(ngapMsg.ProtocolIEs.List); i++ {
				ie := ngapMsg.ProtocolIEs.List[i]
				if ie.Id.Value == ngapType.ProtocolIEIDAMFUENGAPID {
					aMFUENGAPID = ie.Value.AMFUENGAPID
				}
			}

		case ngapType.ProcedureCodePDUSessionResourceModify:
			ngapMsg := successfulOutcome.Value.PDUSessionResourceModifyResponse
			for i := 0; i < len(ngapMsg.ProtocolIEs.List); i++ {
				ie := ngapMsg.ProtocolIEs.List[i]
				if ie.Id.Value == ngapType.ProtocolIEIDAMFUENGAPID {
					aMFUENGAPID = ie.Value.AMFUENGAPID
				}
			}

		case ngapType.ProcedureCodeHandoverResourceAllocation:
			ngapMsg := successfulOutcome.Value.HandoverRequestAcknowledge
			for i := 0; i < len(ngapMsg.ProtocolIEs.List); i++ {
				ie := ngapMsg.ProtocolIEs.List[i]
				if ie.Id.Value == ngapType.ProtocolIEIDAMFUENGAPID {
					aMFUENGAPID = ie.Value.AMFUENGAPID
				}
			}
		}

	case ngapType.NGAPPDUPresentUnsuccessfulOutcome:
		unsuccessfulOutcome := message.UnsuccessfulOutcome
		if unsuccessfulOutcome == nil {
			logger.NgapLog.Errorln("unsuccessfulOutcome is nil")
			return nil
		}
		switch unsuccessfulOutcome.ProcedureCode.Value {
		case ngapType.ProcedureCodeAMFConfigurationUpdate:
		case ngapType.ProcedureCodeInitialContextSetup:
			ngapMsg := unsuccessfulOutcome.Value.InitialContextSetupFailure
			for i := 0; i < len(ngapMsg.ProtocolIEs.List); i++ {
				ie := ngapMsg.ProtocolIEs.List[i]
				if ie.Id.Value == ngapType.ProtocolIEIDAMFUENGAPID {
					aMFUENGAPID = ie.Value.AMFUENGAPID
				}
			}

		case ngapType.ProcedureCodeUEContextModification:
			ngapMsg := unsuccessfulOutcome.Value.UEContextModificationFailure
			for i := 0; i < len(ngapMsg.ProtocolIEs.List); i++ {
				ie := ngapMsg.ProtocolIEs.List[i]
				if ie.Id.Value == ngapType.ProtocolIEIDAMFUENGAPID {
					aMFUENGAPID = ie.Value.AMFUENGAPID
				}
			}

		case ngapType.ProcedureCodeHandoverResourceAllocation:
			ngapMsg := unsuccessfulOutcome.Value.HandoverFailure
			for i := 0; i < len(ngapMsg.ProtocolIEs.List); i++ {
				ie := ngapMsg.ProtocolIEs.List[i]
				if ie.Id.Value == ngapType.ProtocolIEIDAMFUENGAPID {
					aMFUENGAPID = ie.Value.AMFUENGAPID
				}
			}
		}
	}
	return aMFUENGAPID
}
