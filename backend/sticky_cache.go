package backend

import (
	"github.com/omec-project/ngap"
	"github.com/omec-project/ngap/ngapType"
	"github.com/omec-project/sctplb/logger"
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
	logger.AppLog.Infoln("Backend added to cache")
}

func cacheBackend(msg []byte, backend Backend) {
	var ngapId *ngapType.AMFUENGAPID = nil
	ngapMsg, err := ngap.Decoder(msg)
	if err != nil {
		logger.AppLog.Errorln("Cache cannot decode msg: ", err)
		return
	}
	ngapId = extractAMFUENGAPID(ngapMsg)
	addBackend(ngapId, backend)
}

func getCachedBackend(msg []byte) Backend {
	var ngapId *ngapType.AMFUENGAPID = nil
	ngapMsg, err := ngap.Decoder(msg)
	if err != nil {
		logger.AppLog.Errorln("Cache cannot decode msg: ", err)
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

func testMsgDecryption(msg []byte) (int64, bool) {
	ngapMsg, err := ngap.Decoder(msg)
	if err != nil {
		logger.AppLog.Errorln("Cache cannot decode msg: ", err)
		return -1, false
	}
	ngapId := extractAMFUENGAPID(ngapMsg)
	if ngapId == nil {
		return 0, true
	}
	return ngapId.Value, true
}
