package backend

import (
	"github.com/omec-project/ngap"
	"github.com/omec-project/ngap/logger"
	"github.com/omec-project/ngap/ngapType"
	"sync"
)

var (
	cacheMtx sync.Mutex
	cache    map[int64]Backend
)

func init() {
	cache = make(map[int64]Backend)
}

func addBackend(ngapId *ngapType.AMFUENGAPID, backend Backend) {
	if ngapId == nil {
		return
	}

	cacheMtx.Lock()
	defer cacheMtx.Unlock()
	cache[ngapId.Value] = backend
	logger.NgapLog.Infoln("Backend added to cache")
}

func cacheBackend(msg []byte, backend Backend) {
	var ngapId *ngapType.AMFUENGAPID = nil
	ngapMsg, err := ngap.Decoder(msg)
	if err != nil {
		logger.NgapLog.Errorln("Cache cannot decode msg: ", err)
		return
	}
	ngapId = extractAMFUENGAPID(ngapMsg)
	addBackend(ngapId, backend)
}

func getCachedBackend(msg []byte) Backend {
	var ngapId *ngapType.AMFUENGAPID = nil
	ngapMsg, err := ngap.Decoder(msg)
	if err != nil {
		logger.NgapLog.Errorln("Cache cannot decode msg: ", err)
		return nil
	}
	ngapId = extractAMFUENGAPID(ngapMsg)
	if ngapId == nil {
		return nil
	}

	cacheMtx.Lock()
	defer cacheMtx.Unlock()
	return cache[ngapId.Value]
}

// Taken and adapted from ngap/handler.go
func extractAMFUENGAPID(message *ngapType.NGAPPDU) *ngapType.AMFUENGAPID {
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
