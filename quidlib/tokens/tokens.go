package tokens

import (
	emolib "github.com/teal-finance/emo"
)

var emo = emolib.NewZone("tokens")

// Init : init the token zone.
func Init(isVerbose, isDev bool, isCmd bool) {
	if !isDev && !isCmd {
		emo.Print = isVerbose
	}
}
