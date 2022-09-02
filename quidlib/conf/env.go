package conf

import (
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// AdminUser :.
var AdminUser string

// AdminPassword ;.
var AdminPassword string

// InitFromEnv : get the config from environment variables.
func InitFromEnv(isDevMode bool) (conn string, port int) {
	log.Info("Initializing from env")

	if isDevMode {
		log.Fatal("Dev mode is not authorized when initializing from env variables")
	}

	hexKey := os.Getenv("QUID_KEY")
	var err error
	EncodingKey, err = hex.DecodeString(hexKey)
	if err != nil {
		log.Fatal("The key in config must be in hexadecimal format err=", err)
	}

	AdminUser = os.Getenv("QUID_ADMIN_USER")
	AdminPassword = os.Getenv("QUID_ADMIN_PWD")

	conn = os.Getenv("DATABASE_URL")
	conn = strings.Replace(conn, "postgresql://", "postgres://", 1)

	portStr := os.Getenv("PORT")

	port, err = strconv.Atoi(portStr)
	if err != nil {
		log.Panic("PORT must be an integer got=", portStr)
	}

	return conn, port
}
