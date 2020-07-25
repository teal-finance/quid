package tokens

import emolib "github.com/synw/quid/quidlib/emo"

var emo = emolib.Zone{
	Name:    "tokens",
	NoPrint: true,
}

// Init : init the db conf
func Init(isDev bool) {
	emo.NoPrint = !isDev
}
