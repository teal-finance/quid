package api

import (
	emolib "github.com/teal-finance/emo"
)

var emo = emolib.NewZone("api")

// Init : init the db conf.
func Init(isVerbose, isDev bool) {
	if !isDev {
		emo.Print = isVerbose
		emo.Info("Print verbose")
	}
}
