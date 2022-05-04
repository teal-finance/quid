package cmds

import (
	"fmt"
	"os"

	"github.com/teal-finance/quid/quidlib/conf"
)

func GeNConf() {
	fmt.Println("Generating config file")

	if err := conf.Create(); err != nil {
		fmt.Println("Cannot create config file ", err)
		os.Exit(3)
	}

	fmt.Println("Config file created: edit config.json to provide your database settings")
}
