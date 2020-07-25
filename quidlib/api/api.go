package api

import emolib "github.com/synw/quid/quidlib/emo"

var emo = emolib.Zone{
	Name:            "api",
	DeactivatePrint: true,
}

// Init : init the db conf
func Init(isDev bool) {
	emo.DeactivatePrint = !isDev
}
