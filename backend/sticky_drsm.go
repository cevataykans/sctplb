package backend

import (
	"github.com/omec-project/util/drsm"
	"os"
)

var (
	drsmClient drsm.DrsmInterface
)

//func init() {
//	var err error = nil
//	drsmClient, err = initDrsmReadonly()
//	if err != nil {
//		panic(err)
//	}
//}

func initDrsmReadonly() (drsm.DrsmInterface, error) {
	podname := os.Getenv("HOSTNAME")
	podip := os.Getenv("POD_IP")

	// The LB doesn't need a unique NFID like AMF does
	lbPodId := drsm.PodId{
		PodName:     podname,
		PodInstance: "sctplb-load-balancer",
		PodIp:       podip,
	}

	dbUrl := "mongodb://mongodb-arbiter-headless"
	db := drsm.DbInfo{
		Url:  dbUrl,
		Name: "sdcore_amf",
	}

	// Use Demux mode (read-only)
	opt := &drsm.Options{
		ResIdSize: 24,
		Mode:      drsm.ResourceDemux,
	}

	return drsm.InitDRSM("amfid", lbPodId, db, opt)
}
