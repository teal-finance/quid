package conf

import (
	"encoding/hex"
	"fmt"
	"log"
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
	fmt.Println("Initializing from env")

	if isDevMode {
		fmt.Println("Dev mode is not authorized when initializing from env variables")
		os.Exit(1)
	}

	hexKey := os.Getenv("QUID_KEY")
	var err error
	EncodingKey, err = hex.DecodeString(hexKey)
	if err != nil {
		fmt.Println("The key in config must be in hexadecimal format err=", err)
		os.Exit(5)
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
