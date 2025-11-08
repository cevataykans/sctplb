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
}

func cacheBackend(msg []byte, backend Backend) {
	var ngapId *ngapType.AMFUENGAPID = nil
	ngapMsg, err := ngap.Decoder(msg)
	if err != nil {
		return
	}
	ngapId = extractAMFUENGAPID(ngapMsg)
	addBackend(ngapId, backend)
	logger.AppLog.Infoln("Backend added to cache")
}

func getCachedBackend(msg []byte) Backend {
	var ngapId *ngapType.AMFUENGAPID = nil
	ngapMsg, err := ngap.Decoder(msg)
	if err != nil {
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
