package api

import (
	"github.com/teal-finance/emo"
)

var logg = emo.NewLogger("api")

// Init : init the db conf.
func Init(isVerbose, isDev bool) {
	if !isDev {
		logg.Print = isVerbose
		logg.Info("Print verbose")
	}
}
