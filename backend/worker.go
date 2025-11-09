package backend

import (
	"github.com/omec-project/ngap"
	"github.com/omec-project/ngap/ngapType"
	"github.com/omec-project/sctplb/context"
	"github.com/omec-project/sctplb/logger"
)

var (
	ranMsgs chan *task
)

func init() {
	ranMsgs = make(chan *task, 1024)
}

type task struct {
	msg []byte
	ran *context.Ran
}

func worker() {
	for task := range ranMsgs {
		msg := task.msg
		ran := task.ran

		var ngapId *ngapType.AMFUENGAPID = nil
		ngapMsg, err := ngap.Decoder(msg)
		if err == nil {
			ngapId = extractAMFUENGAPID(ngapMsg)
		}

		drsmBackend, err := findBackendWithNGAPID(ngapId)
		if err == nil {
			// send msg to the returned backend
			err = drsmBackend.Send(msg, false, ran)
			if err != nil {
				logger.SctpLog.Errorln("can not send to backend returned by drsm:", err)
			}
			return
		}

		for {
			// Select the backend NF based on RoundRobin Algorithm
			backend := RoundRobin()
			if backend == nil {
				break
			}

			if backend.State() {
				if err := backend.Send(msg, false, ran); err != nil {
					logger.SctpLog.Errorln("can not send:", err)
				}
				break
			}
		}
	}
}
