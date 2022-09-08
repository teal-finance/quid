package api

import (
	"github.com/teal-finance/emo"
)

var log = emo.NewZone("api")

// Init : init the db conf.
func Init(isVerbose, isDev bool) {
	if !isDev {
		log.Verbose = emo.No
		log.Info("Print verbose")
	}
}
