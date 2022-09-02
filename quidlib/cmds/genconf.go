package cmds

import (
	"github.com/teal-finance/quid/quidlib/conf"
)

func GeNConf() {
	log.Info("Generating config file")

	if err := conf.Create(); err != nil {
		log.Fatal("Cannot create config file", err)
	}

	log.State("Config file created: edit config.json to provide your database settings")
}
