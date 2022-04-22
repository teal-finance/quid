package tokens

import (
	emolib "github.com/teal-finance/emo"
)

var Emo = emolib.NewZone("tokens")

// Init : init the db conf.
func Init(isVerbose bool, isDev bool) {
	if !isDev {
		Emo.Print = isVerbose
	}
}
