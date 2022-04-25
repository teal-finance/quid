package api

import (
	emolib "github.com/teal-finance/emo"
)

var emo = emolib.NewZone("api")

// Init : init the db conf.
func Init(isVerbose bool, isDev bool) {
	if !isDev {
		emo.Print = isVerbose
	}
}
