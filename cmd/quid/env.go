package main

import (
	"encoding/hex"
	"os"
	"strconv"
	"strings"

	"github.com/teal-finance/quid/crypt"
)

// AdminUser :.
var adminUser string

// AdminPassword ;.
var adminPassword string

// InitFromEnv : get the config from environment variables.
func initFromEnv(isDevMode bool) (conn string, port int) {
	log.Info("Initializing from env")

	if isDevMode {
		log.Fatal("Dev mode is not authorized when initializing from env variables")
	}

	hexKey := os.Getenv("QUID_KEY")
	if len(hexKey) < 32 {
		log.Panic("Want AES-128 key composed by 32 hexadecimal digits, but got", len(hexKey))
	}

	var err error
	crypt.EncodingKey, err = hex.DecodeString(hexKey[:32])
	if err != nil {
		log.Fatal("The key in config must be in hexadecimal format err=", err)
	}

	adminUser = os.Getenv("QUID_ADMIN_USER")
	adminPassword = os.Getenv("QUID_ADMIN_PWD")

	conn = os.Getenv("DATABASE_URL")
	conn = strings.Replace(conn, "postgresql://", "postgres://", 1)

	portStr := os.Getenv("PORT")

	port, err = strconv.Atoi(portStr)
	if err != nil {
		log.Panic("PORT must be an integer got=", portStr)
	}

	return conn, port
}
