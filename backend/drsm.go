package backend

import (
	"fmt"
	"github.com/omec-project/ngap"
	"github.com/omec-project/ngap/logger"
	"github.com/omec-project/ngap/ngapType"
	"github.com/omec-project/sctplb/context"
	"sync"
	//"github.com/omec-project/util/drsm"
)

var (
	//drsmClient drsm.DrsmInterface

	initiatingMsgParsers = map[int64]func(val *ngapType.InitiatingMessageValue) *ngapType.AMFUENGAPID{
		ngapType.ProcedureCodeUplinkNASTransport:                 initialUeMsgUplinkNASTransport,
		ngapType.ProcedureCodeHandoverCancel:                     initialUeMsgHandoverCancel,
		ngapType.ProcedureCodeUEContextReleaseRequest:            initialUeMsgUEContextReleaseRequest,
		ngapType.ProcedureCodeNASNonDeliveryIndication:           initialUeMsgNASNonDeliveryIndication,
		ngapType.ProcedureCodeUERadioCapabilityInfoIndication:    initialUeMsgUERadioCapabilityInfoIndication,
		ngapType.ProcedureCodeHandoverNotification:               initialUeMsgHandoverNotify,
		ngapType.ProcedureCodeHandoverPreparation:                initialUeMsgHandoverRequired,
		ngapType.ProcedureCodePDUSessionResourceNotify:           initialUeMsgPDUSessionResourceNotify,
		ngapType.ProcedureCodePathSwitchRequest:                  initialUeMsgPathSwitchRequest,
		ngapType.ProcedureCodePDUSessionResourceModifyIndication: initialUeMsgPDUSessionResourceModifyIndication,
	}

	successMsgParsers = map[int64]func(val *ngapType.SuccessfulOutcomeValue) *ngapType.AMFUENGAPID{
		ngapType.ProcedureCodeInitialContextSetup:        successMsgInitialContextSetupResponse,
		ngapType.ProcedureCodeUEContextModification:      successMsgUEContextModificationResponse,
		ngapType.ProcedureCodePDUSessionResourceSetup:    successMsgPDUSessionResourceSetupResponse,
		ngapType.ProcedureCodePDUSessionResourceModify:   successMsgPDUSessionResourceModifyResponse,
		ngapType.ProcedureCodeHandoverResourceAllocation: successMsgHandoverRequestAcknowledge,
		ngapType.ProcedureCodeUEContextRelease:           successMsgUEContextReleaseComplete,
		ngapType.ProcedureCodePDUSessionResourceRelease:  successMsgPDUSessionResourceReleaseResponse,
	}

	unsuccessMsgParsers = map[int64]func(val *ngapType.UnsuccessfulOutcomeValue) *ngapType.AMFUENGAPID{
		ngapType.ProcedureCodeInitialContextSetup:        unsuccessMsgInitialContextSetupFailure,
		ngapType.ProcedureCodeUEContextModification:      unsuccessMsgUEContextModificationFailure,
		ngapType.ProcedureCodeHandoverResourceAllocation: unsuccessMsgHandoverFailure,
	}

	cacheLock = sync.RWMutex{}
	cache     = make(map[int32]string)
)

//func init() {
//	var err error = nil
//	drsmClient, err = initDrsmReadonly()
//	if err != nil {
//		panic(err)
//	}
//}

func tryCacheMsg(msg []byte, idPodIp string) {
	var ngapId *ngapType.AMFUENGAPID = nil
	ngapMsg, err := ngap.Decoder(msg)
	if err != nil {
		return
	}
	ngapId = extractAMFUENGAPID(ngapMsg)
	if ngapId != nil {
		placeInCache(int32(ngapId.Value), idPodIp)
	}
}

func placeInCache(ngapIDVal int32, idPodIp string) {
	cacheLock.Lock()
	cache[ngapIDVal] = idPodIp
	cacheLock.Unlock()
}

func findInCache(ngapIDVal int32) string {
	cacheLock.RLock()
	id := cache[ngapIDVal]
	cacheLock.RUnlock()
	return id
}

func findBackendWithNGAPID(ctx *context.SctplbContext, ngapId *ngapType.AMFUENGAPID) (Backend, error) {
	if ngapId == nil {
		return nil, fmt.Errorf("ngapId is nil")
	}

	//id, err := drsmClient.FindOwnerInt32ID(int32(ngapId.Value))
	idPodIp := findInCache(int32(ngapId.Value))
	if idPodIp == "" {
		return nil, fmt.Errorf("ngapId not found in cache")
	}
	logger.NgapLog.Infoln("Found backend with id:", idPodIp)

	for _, instance := range ctx.Backends {
		b1 := instance.(*GrpcServer)
		// AMF sets RedirectID as PodIp
		if b1.address != idPodIp {
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

//func initDrsmReadonly() (drsm.DrsmInterface, error) {
//	podname := os.Getenv("HOSTNAME")
//	podip := os.Getenv("POD_IP")
//
//	// The LB doesn't need a unique NFID like AMF does
//	lbPodId := drsm.PodId{
//		PodName:     podname,
//		PodInstance: "sctplb-load-balancer",
//		PodIp:       podip,
//	}
//
//	dbUrl := "mongodb://mongodb-arbiter-headless"
//	db := drsm.DbInfo{
//		Url:  dbUrl,
//		Name: "sdcore_amf",
//	}
//
//	// Use Demux mode (read-only)
//	opt := &drsm.Options{
//		ResIdSize: 24,
//		Mode:      drsm.ResourceDemux,
//	}
//
//	return drsm.InitDRSM("amfid", lbPodId, db, opt)
//}

// Taken and adapted from ngap/handler.go
func extractAMFUENGAPID(message *ngapType.NGAPPDU) *ngapType.AMFUENGAPID {
	if message == nil {
		logger.NgapLog.Errorln("NGAP Message is nil")
		return nil
	}

	var aMFUENGAPID *ngapType.AMFUENGAPID = nil
	switch message.Present {
	case ngapType.NGAPPDUPresentInitiatingMessage:
		initiatingMessage := message.InitiatingMessage
		if initiatingMessage == nil {
			logger.NgapLog.Errorln("initiatingMessage is nil")
			return nil
		}
		parser := initiatingMsgParsers[initiatingMessage.ProcedureCode.Value]
		if parser != nil {
			return parser(&initiatingMessage.Value)
		}
	case ngapType.NGAPPDUPresentSuccessfulOutcome:
		successfulOutcome := message.SuccessfulOutcome
		if successfulOutcome == nil {
			logger.NgapLog.Errorln("successfulOutcome is nil")
			return nil
		}
		parser := successMsgParsers[successfulOutcome.ProcedureCode.Value]
		if parser != nil {
			return parser(&successfulOutcome.Value)
		}
	case ngapType.NGAPPDUPresentUnsuccessfulOutcome:
		unsuccessfulOutcome := message.UnsuccessfulOutcome
		if unsuccessfulOutcome == nil {
			logger.NgapLog.Errorln("unsuccessfulOutcome is nil")
			return nil
		}
		parser := unsuccessMsgParsers[unsuccessfulOutcome.ProcedureCode.Value]
		if parser != nil {
			return parser(&unsuccessfulOutcome.Value)
		}
	}
	return aMFUENGAPID
}

func initialUeMsgUplinkNASTransport(val *ngapType.InitiatingMessageValue) *ngapType.AMFUENGAPID {
	msg := val.UplinkNASTransport
	for _, ie := range msg.ProtocolIEs.List {
		if ie.Id.Value == ngapType.ProtocolIEIDAMFUENGAPID {
			return ie.Value.AMFUENGAPID
		}
	}
	return nil
}

func initialUeMsgHandoverCancel(val *ngapType.InitiatingMessageValue) *ngapType.AMFUENGAPID {
	msg := val.HandoverCancel
	for _, ie := range msg.ProtocolIEs.List {
		if ie.Id.Value == ngapType.ProtocolIEIDAMFUENGAPID {
			return ie.Value.AMFUENGAPID
		}
	}
	return nil
}

func initialUeMsgUEContextReleaseRequest(val *ngapType.InitiatingMessageValue) *ngapType.AMFUENGAPID {
	msg := val.UEContextReleaseRequest
	for _, ie := range msg.ProtocolIEs.List {
		if ie.Id.Value == ngapType.ProtocolIEIDAMFUENGAPID {
			return ie.Value.AMFUENGAPID
		}
	}
	return nil
}

func initialUeMsgNASNonDeliveryIndication(val *ngapType.InitiatingMessageValue) *ngapType.AMFUENGAPID {
	msg := val.NASNonDeliveryIndication
	for _, ie := range msg.ProtocolIEs.List {
		if ie.Id.Value == ngapType.ProtocolIEIDAMFUENGAPID {
			return ie.Value.AMFUENGAPID
		}
	}
	return nil
}

func initialUeMsgUERadioCapabilityInfoIndication(val *ngapType.InitiatingMessageValue) *ngapType.AMFUENGAPID {
	msg := val.UERadioCapabilityInfoIndication
	for _, ie := range msg.ProtocolIEs.List {
		if ie.Id.Value == ngapType.ProtocolIEIDAMFUENGAPID {
			return ie.Value.AMFUENGAPID
		}
	}
	return nil
}

func initialUeMsgHandoverNotify(val *ngapType.InitiatingMessageValue) *ngapType.AMFUENGAPID {
	msg := val.HandoverNotify
	for _, ie := range msg.ProtocolIEs.List {
		if ie.Id.Value == ngapType.ProtocolIEIDAMFUENGAPID {
			return ie.Value.AMFUENGAPID
		}
	}
	return nil
}

func initialUeMsgHandoverRequired(val *ngapType.InitiatingMessageValue) *ngapType.AMFUENGAPID {
	msg := val.HandoverRequired
	for _, ie := range msg.ProtocolIEs.List {
		if ie.Id.Value == ngapType.ProtocolIEIDAMFUENGAPID {
			return ie.Value.AMFUENGAPID
		}
	}
	return nil
}

func initialUeMsgPDUSessionResourceNotify(val *ngapType.InitiatingMessageValue) *ngapType.AMFUENGAPID {
	msg := val.PDUSessionResourceNotify
	for _, ie := range msg.ProtocolIEs.List {
		if ie.Id.Value == ngapType.ProtocolIEIDAMFUENGAPID {
			return ie.Value.AMFUENGAPID
		}
	}
	return nil
}

func initialUeMsgPathSwitchRequest(val *ngapType.InitiatingMessageValue) *ngapType.AMFUENGAPID {
	msg := val.PathSwitchRequest
	for _, ie := range msg.ProtocolIEs.List {
		if ie.Id.Value == ngapType.ProtocolIEIDSourceAMFUENGAPID {
			return ie.Value.SourceAMFUENGAPID
		}
	}
	return nil
}

func initialUeMsgPDUSessionResourceModifyIndication(val *ngapType.InitiatingMessageValue) *ngapType.AMFUENGAPID {
	msg := val.PDUSessionResourceModifyIndication
	for _, ie := range msg.ProtocolIEs.List {
		if ie.Id.Value == ngapType.ProtocolIEIDAMFUENGAPID {
			return ie.Value.AMFUENGAPID
		}
	}
	return nil
}

func successMsgUEContextReleaseComplete(val *ngapType.SuccessfulOutcomeValue) *ngapType.AMFUENGAPID {
	msg := val.UEContextReleaseComplete
	for _, ie := range msg.ProtocolIEs.List {
		if ie.Id.Value == ngapType.ProtocolIEIDAMFUENGAPID {
			return ie.Value.AMFUENGAPID
		}
	}
	return nil
}

func successMsgPDUSessionResourceReleaseResponse(val *ngapType.SuccessfulOutcomeValue) *ngapType.AMFUENGAPID {
	msg := val.PDUSessionResourceReleaseResponse
	for _, ie := range msg.ProtocolIEs.List {
		if ie.Id.Value == ngapType.ProtocolIEIDAMFUENGAPID {
			return ie.Value.AMFUENGAPID
		}
	}
	return nil
}

func successMsgInitialContextSetupResponse(val *ngapType.SuccessfulOutcomeValue) *ngapType.AMFUENGAPID {
	msg := val.InitialContextSetupResponse
	for _, ie := range msg.ProtocolIEs.List {
		if ie.Id.Value == ngapType.ProtocolIEIDAMFUENGAPID {
			return ie.Value.AMFUENGAPID
		}
	}
	return nil
}

func successMsgUEContextModificationResponse(val *ngapType.SuccessfulOutcomeValue) *ngapType.AMFUENGAPID {
	msg := val.UEContextModificationResponse
	for _, ie := range msg.ProtocolIEs.List {
		if ie.Id.Value == ngapType.ProtocolIEIDAMFUENGAPID {
			return ie.Value.AMFUENGAPID
		}
	}
	return nil
}

func successMsgPDUSessionResourceSetupResponse(val *ngapType.SuccessfulOutcomeValue) *ngapType.AMFUENGAPID {
	msg := val.PDUSessionResourceSetupResponse
	for _, ie := range msg.ProtocolIEs.List {
		if ie.Id.Value == ngapType.ProtocolIEIDAMFUENGAPID {
			return ie.Value.AMFUENGAPID
		}
	}
	return nil
}

func successMsgPDUSessionResourceModifyResponse(val *ngapType.SuccessfulOutcomeValue) *ngapType.AMFUENGAPID {
	msg := val.PDUSessionResourceModifyResponse
	for _, ie := range msg.ProtocolIEs.List {
		if ie.Id.Value == ngapType.ProtocolIEIDAMFUENGAPID {
			return ie.Value.AMFUENGAPID
		}
	}
	return nil
}

func successMsgHandoverRequestAcknowledge(val *ngapType.SuccessfulOutcomeValue) *ngapType.AMFUENGAPID {
	msg := val.HandoverRequestAcknowledge
	for _, ie := range msg.ProtocolIEs.List {
		if ie.Id.Value == ngapType.ProtocolIEIDAMFUENGAPID {
			return ie.Value.AMFUENGAPID
		}
	}
	return nil
}

func unsuccessMsgInitialContextSetupFailure(val *ngapType.UnsuccessfulOutcomeValue) *ngapType.AMFUENGAPID {
	msg := val.InitialContextSetupFailure
	if msg == nil {
		logger.NgapLog.Errorln("InitialContextSetupFailure is nil")
		return nil
	}
	for _, ie := range msg.ProtocolIEs.List {
		if ie.Id.Value == ngapType.ProtocolIEIDAMFUENGAPID {
			return ie.Value.AMFUENGAPID
		}
	}
	return nil
}

func unsuccessMsgUEContextModificationFailure(val *ngapType.UnsuccessfulOutcomeValue) *ngapType.AMFUENGAPID {
	msg := val.UEContextModificationFailure
	if msg == nil {
		logger.NgapLog.Errorln("UEContextModificationFailure is nil")
		return nil
	}
	for _, ie := range msg.ProtocolIEs.List {
		if ie.Id.Value == ngapType.ProtocolIEIDAMFUENGAPID {
			return ie.Value.AMFUENGAPID
		}
	}
	return nil
}

func unsuccessMsgHandoverFailure(val *ngapType.UnsuccessfulOutcomeValue) *ngapType.AMFUENGAPID {
	msg := val.HandoverFailure
	if msg == nil {
		logger.NgapLog.Errorln("HandoverFailure is nil")
		return nil
	}
	for _, ie := range msg.ProtocolIEs.List {
		if ie.Id.Value == ngapType.ProtocolIEIDAMFUENGAPID {
			return ie.Value.AMFUENGAPID
		}
	}
	return nil
}
