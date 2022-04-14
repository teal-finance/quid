package tokens

import emolib "github.com/teal-finance/quid/quidlib/emo"

var emo = emolib.NewZone("tokens")

// Init : init the db conf.
func Init(isVerbose bool, isDev bool) {
	if !isDev {
		emo.Print = isVerbose
	}
}
