package conf

import (
	"fmt"
	"os"
	"strings"
)

// AdminUser :.
var AdminUser string

// AdminPassword ;.
var AdminPassword string

// InitFromEnv : get the config from environment variables.
func InitFromEnv(isDevMode bool) (conn, port string) {
	fmt.Println("Initializing from env")

	if isDevMode {
		fmt.Println("Dev mode is not authorized when initializing from env variables")
		os.Exit(1)
	}

	EncodingKey = os.Getenv("QUID_KEY")

	AdminUser = os.Getenv("QUID_ADMIN_USER")
	AdminPassword = os.Getenv("QUID_ADMIN_PWD")

	conn = os.Getenv("DATABASE_URL")
	conn = strings.Replace(conn, "postgresql://", "postgres://", 1)

	port = os.Getenv("PORT")

	return conn, port
}
