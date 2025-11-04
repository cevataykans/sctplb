package backend

import (
	"github.com/omec-project/ngap/logger"
	"github.com/omec-project/ngap/ngapType"
)

//type stickySession struct {
//	ueId  int64
//	ranId int64
//}

var (
	stickySessions = make(map[string]Backend)
)

func extractUEIdentifier(ranMsg *ngapType.NGAPPDU) *ngapType.RANUENGAPID {
	var rANUENGAPID *ngapType.RANUENGAPID

	switch ranMsg.Present {
	case ngapType.NGAPPDUPresentInitiatingMessage:
		initiatingMessage := ranMsg.InitiatingMessage
		if initiatingMessage == nil {
			return nil
		}
		switch initiatingMessage.ProcedureCode.Value {
		case ngapType.ProcedureCodeNGSetup:
		case ngapType.ProcedureCodeInitialUEMessage:
			ngapMsg := initiatingMessage.Value.InitialUEMessage
			if ngapMsg == nil {
				logger.NgapLog.Errorln("initialUEMessage is nil")
				return nil
			}
			for i := 0; i < len(ngapMsg.ProtocolIEs.List); i++ {
				ie := ngapMsg.ProtocolIEs.List[i]
				if ie.Id.Value == ngapType.ProtocolIEIDRANUENGAPID {
					rANUENGAPID = ie.Value.RANUENGAPID
					if rANUENGAPID == nil {
						logger.NgapLog.Errorln("ranUeNgapID is nil")
						return nil
					}
				}
			}
		case ngapType.ProcedureCodeUplinkNASTransport:
			ngapMsg := initiatingMessage.Value.UplinkNASTransport
			if ngapMsg == nil {
				logger.NgapLog.Errorln("UplinkNasTransport is nil")
				return nil
			}
			for i := 0; i < len(ngapMsg.ProtocolIEs.List); i++ {
				ie := ngapMsg.ProtocolIEs.List[i]
				if ie.Id.Value == ngapType.ProtocolIEIDRANUENGAPID {
					rANUENGAPID = ie.Value.RANUENGAPID
					logger.NgapLog.Debugln("decode IE RanUeNgapID")
					if rANUENGAPID == nil {
						logger.NgapLog.Errorln("RanUeNgapID is nil")
						return nil
					}
				}
			}
		case ngapType.ProcedureCodeHandoverCancel:
			ngapMsg := initiatingMessage.Value.HandoverCancel
			for i := 0; i < len(ngapMsg.ProtocolIEs.List); i++ {
				ie := ngapMsg.ProtocolIEs.List[i]
				if ie.Id.Value == ngapType.ProtocolIEIDRANUENGAPID {
					rANUENGAPID = ie.Value.RANUENGAPID
					logger.NgapLog.Debugln("decode IE RanUeNgapID")
					if rANUENGAPID == nil {
						logger.NgapLog.Errorln("RANUENGAPID is nil")
						return nil
					}
				}
			}
		case ngapType.ProcedureCodeUEContextReleaseRequest:
			ngapMsg := initiatingMessage.Value.UEContextReleaseRequest
			for i := 0; i < len(ngapMsg.ProtocolIEs.List); i++ {
				ie := ngapMsg.ProtocolIEs.List[i]
				if ie.Id.Value == ngapType.ProtocolIEIDRANUENGAPID {
					rANUENGAPID = ie.Value.RANUENGAPID
					logger.NgapLog.Debugln("decode IE RanUeNgapID")
					if rANUENGAPID == nil {
						logger.NgapLog.Errorln("RANUENGAPID is nil")
						return nil
					}
				}
			}
		case ngapType.ProcedureCodeNASNonDeliveryIndication:
			ngapMsg := initiatingMessage.Value.NASNonDeliveryIndication
			for i := 0; i < len(ngapMsg.ProtocolIEs.List); i++ {
				ie := ngapMsg.ProtocolIEs.List[i]
				if ie.Id.Value == ngapType.ProtocolIEIDRANUENGAPID {
					rANUENGAPID = ie.Value.RANUENGAPID
					logger.NgapLog.Debugln("decode IE RanUeNgapID")
					if rANUENGAPID == nil {
						logger.NgapLog.Errorln("RANUENGAPID is nil")
						return nil
					}
				}
			}
		case ngapType.ProcedureCodeLocationReportingFailureIndication:
		case ngapType.ProcedureCodeErrorIndication:
		case ngapType.ProcedureCodeUERadioCapabilityInfoIndication:
			ngapMsg := initiatingMessage.Value.UERadioCapabilityInfoIndication
			for i := 0; i < len(ngapMsg.ProtocolIEs.List); i++ {
				ie := ngapMsg.ProtocolIEs.List[i]
				if ie.Id.Value == ngapType.ProtocolIEIDRANUENGAPID {
					rANUENGAPID = ie.Value.RANUENGAPID
					logger.NgapLog.Debugln("decode IE RanUeNgapID")
					if rANUENGAPID == nil {
						logger.NgapLog.Errorln("RANUENGAPID is nil")
						return nil
					}
				}
			}
		case ngapType.ProcedureCodeHandoverNotification:
			ngapMsg := initiatingMessage.Value.HandoverNotify
			for i := 0; i < len(ngapMsg.ProtocolIEs.List); i++ {
				ie := ngapMsg.ProtocolIEs.List[i]
				if ie.Id.Value == ngapType.ProtocolIEIDRANUENGAPID {
					rANUENGAPID = ie.Value.RANUENGAPID
					logger.NgapLog.Debugln("decode IE RanUeNgapID")
					if rANUENGAPID == nil {
						logger.NgapLog.Errorln("RANUENGAPID is nil")
						return nil
					}
				}
			}
		case ngapType.ProcedureCodeHandoverPreparation:
			ngapMsg := initiatingMessage.Value.HandoverRequired
			for i := 0; i < len(ngapMsg.ProtocolIEs.List); i++ {
				ie := ngapMsg.ProtocolIEs.List[i]
				if ie.Id.Value == ngapType.ProtocolIEIDRANUENGAPID {
					rANUENGAPID = ie.Value.RANUENGAPID
					logger.NgapLog.Debugln("decode IE RanUeNgapID")
					if rANUENGAPID == nil {
						logger.NgapLog.Errorln("RANUENGAPID is nil")
						return nil
					}
				}
			}
		case ngapType.ProcedureCodeRANConfigurationUpdate:
		case ngapType.ProcedureCodeRRCInactiveTransitionReport:
		case ngapType.ProcedureCodePDUSessionResourceNotify:
			ngapMsg := initiatingMessage.Value.PDUSessionResourceNotify
			for i := 0; i < len(ngapMsg.ProtocolIEs.List); i++ {
				ie := ngapMsg.ProtocolIEs.List[i]
				if ie.Id.Value == ngapType.ProtocolIEIDRANUENGAPID {
					rANUENGAPID = ie.Value.RANUENGAPID
					logger.NgapLog.Debugln("decode IE RanUeNgapID")
					if rANUENGAPID == nil {
						logger.NgapLog.Errorln("RANUENGAPID is nil")
						return nil
					}
				}
			}
		case ngapType.ProcedureCodePathSwitchRequest:
			ngapMsg := initiatingMessage.Value.PathSwitchRequest
			for i := 0; i < len(ngapMsg.ProtocolIEs.List); i++ {
				ie := ngapMsg.ProtocolIEs.List[i]
				if ie.Id.Value == ngapType.ProtocolIEIDRANUENGAPID {
					rANUENGAPID = ie.Value.RANUENGAPID
					logger.NgapLog.Debugln("decode IE RanUeNgapID")
					if rANUENGAPID == nil {
						logger.NgapLog.Errorln("RANUENGAPID is nil")
						return nil
					}
				}
			}
		case ngapType.ProcedureCodeLocationReport:
		case ngapType.ProcedureCodeUplinkUEAssociatedNRPPaTransport:
		case ngapType.ProcedureCodeUplinkRANConfigurationTransfer:
		case ngapType.ProcedureCodePDUSessionResourceModifyIndication:
			ngapMsg := initiatingMessage.Value.PDUSessionResourceModifyIndication
			for i := 0; i < len(ngapMsg.ProtocolIEs.List); i++ {
				ie := ngapMsg.ProtocolIEs.List[i]
				if ie.Id.Value == ngapType.ProtocolIEIDRANUENGAPID {
					rANUENGAPID = ie.Value.RANUENGAPID
					logger.NgapLog.Debugln("decode IE RanUeNgapID")
					if rANUENGAPID == nil {
						logger.NgapLog.Errorln("RANUENGAPID is nil")
						return nil
					}
				}
			}
		case ngapType.ProcedureCodeCellTrafficTrace:
		case ngapType.ProcedureCodeUplinkRANStatusTransfer:
		case ngapType.ProcedureCodeUplinkNonUEAssociatedNRPPaTransport:
		}

	case ngapType.NGAPPDUPresentSuccessfulOutcome:
		successfulOutcome := ranMsg.SuccessfulOutcome
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
				if ie.Id.Value == ngapType.ProtocolIEIDRANUENGAPID {
					rANUENGAPID = ie.Value.RANUENGAPID
					logger.NgapLog.Debugln("decode IE RanUeNgapID")
					if rANUENGAPID == nil {
						logger.NgapLog.Errorln("RANUENGAPID is nil")
						return nil
					}
				}
			}

		case ngapType.ProcedureCodePDUSessionResourceRelease:
			ngapMsg := successfulOutcome.Value.PDUSessionResourceReleaseResponse
			for i := 0; i < len(ngapMsg.ProtocolIEs.List); i++ {
				ie := ngapMsg.ProtocolIEs.List[i]
				if ie.Id.Value == ngapType.ProtocolIEIDRANUENGAPID {
					rANUENGAPID = ie.Value.RANUENGAPID
					logger.NgapLog.Debugln("decode IE RanUeNgapID")
					if rANUENGAPID == nil {
						logger.NgapLog.Errorln("RANUENGAPID is nil")
						return nil
					}
				}
			}

		case ngapType.ProcedureCodeUERadioCapabilityCheck:
		case ngapType.ProcedureCodeAMFConfigurationUpdate:
		case ngapType.ProcedureCodeInitialContextSetup:
			ngapMsg := successfulOutcome.Value.InitialContextSetupResponse
			for i := 0; i < len(ngapMsg.ProtocolIEs.List); i++ {
				ie := ngapMsg.ProtocolIEs.List[i]
				if ie.Id.Value == ngapType.ProtocolIEIDRANUENGAPID {
					rANUENGAPID = ie.Value.RANUENGAPID
					logger.NgapLog.Debugln("decode IE RanUeNgapID")
					if rANUENGAPID == nil {
						logger.NgapLog.Errorln("RANUENGAPID is nil")
						return nil
					}
				}
			}

		case ngapType.ProcedureCodeUEContextModification:
			ngapMsg := successfulOutcome.Value.UEContextModificationResponse
			for i := 0; i < len(ngapMsg.ProtocolIEs.List); i++ {
				ie := ngapMsg.ProtocolIEs.List[i]
				if ie.Id.Value == ngapType.ProtocolIEIDRANUENGAPID {
					rANUENGAPID = ie.Value.RANUENGAPID
					logger.NgapLog.Debugln("decode IE RanUeNgapID")
					if rANUENGAPID == nil {
						logger.NgapLog.Errorln("RANUENGAPID is nil")
						return nil
					}
				}
			}

		case ngapType.ProcedureCodePDUSessionResourceSetup:
			ngapMsg := successfulOutcome.Value.PDUSessionResourceSetupResponse
			for i := 0; i < len(ngapMsg.ProtocolIEs.List); i++ {
				ie := ngapMsg.ProtocolIEs.List[i]
				if ie.Id.Value == ngapType.ProtocolIEIDRANUENGAPID {
					rANUENGAPID = ie.Value.RANUENGAPID
					logger.NgapLog.Debugln("decode IE RanUeNgapID")
					if rANUENGAPID == nil {
						logger.NgapLog.Errorln("RANUENGAPID is nil")
						return nil
					}
				}
			}

		case ngapType.ProcedureCodePDUSessionResourceModify:
			ngapMsg := successfulOutcome.Value.PDUSessionResourceModifyResponse
			for i := 0; i < len(ngapMsg.ProtocolIEs.List); i++ {
				ie := ngapMsg.ProtocolIEs.List[i]
				if ie.Id.Value == ngapType.ProtocolIEIDRANUENGAPID {
					rANUENGAPID = ie.Value.RANUENGAPID
					logger.NgapLog.Debugln("decode IE RanUeNgapID")
					if rANUENGAPID == nil {
						logger.NgapLog.Errorln("RANUENGAPID is nil")
						return nil
					}
				}
			}

		case ngapType.ProcedureCodeHandoverResourceAllocation:
			ngapMsg := successfulOutcome.Value.HandoverRequestAcknowledge
			for i := 0; i < len(ngapMsg.ProtocolIEs.List); i++ {
				ie := ngapMsg.ProtocolIEs.List[i]
				if ie.Id.Value == ngapType.ProtocolIEIDRANUENGAPID {
					rANUENGAPID = ie.Value.RANUENGAPID
					logger.NgapLog.Debugln("decode IE RanUeNgapID")
					if rANUENGAPID == nil {
						logger.NgapLog.Errorln("RANUENGAPID is nil")
						return nil
					}
				}
			}
		}
	case ngapType.NGAPPDUPresentUnsuccessfulOutcome:
		logger.NgapLog.Infoln("ngapType.NGAPPDUPresentUnsuccessfulOutcome received for registration procedure")
		unsuccessfulOutcome := ranMsg.UnsuccessfulOutcome
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
				if ie.Id.Value == ngapType.ProtocolIEIDRANUENGAPID {
					rANUENGAPID = ie.Value.RANUENGAPID
					logger.NgapLog.Debugln("decode IE RanUeNgapID")
					if rANUENGAPID == nil {
						logger.NgapLog.Errorln("RANUENGAPID is nil")
						return nil
					}
				}
			}

		case ngapType.ProcedureCodeUEContextModification:
			ngapMsg := unsuccessfulOutcome.Value.UEContextModificationFailure
			for i := 0; i < len(ngapMsg.ProtocolIEs.List); i++ {
				ie := ngapMsg.ProtocolIEs.List[i]
				if ie.Id.Value == ngapType.ProtocolIEIDRANUENGAPID {
					rANUENGAPID = ie.Value.RANUENGAPID
					logger.NgapLog.Debugln("decode IE RanUeNgapID")
					if rANUENGAPID == nil {
						logger.NgapLog.Errorln("RANUENGAPID is nil")
						return nil
					}
				}
			}

		case ngapType.ProcedureCodeHandoverResourceAllocation:
		}
	}
	return rANUENGAPID
}
