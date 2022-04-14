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
func InitFromEnv() (conn, port string) {
	fmt.Println("Initializing from env")

	EncodingKey = os.Getenv("QUID_KEY")

	AdminUser = os.Getenv("QUID_ADMIN_USER")
	AdminPassword = os.Getenv("QUID_ADMIN_PWD")

	conn = os.Getenv("DATABASE_URL")
	conn = strings.Replace(conn, "postgresql://", "", 1)

	port = os.Getenv("PORT")

	return conn, port
}
