package tokens

import (
	"github.com/teal-finance/emo"
)

var logg = emo.NewLogger("tokens")

// Init : init the token zone.
func Init(isVerbose, isDev, isCmd bool) {
	if !isDev && !isCmd {
		logg.Print = isVerbose
	}
}
