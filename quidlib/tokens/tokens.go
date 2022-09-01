package tokens

import (
	"github.com/teal-finance/emo"
)

var log = emo.NewZone("tkn")

// Init : init the token zone.
func Init(isVerbose, isDev, isCmd bool) {
	if !isDev && !isCmd {
		log.Verbose = emo.Yes
	}
}
