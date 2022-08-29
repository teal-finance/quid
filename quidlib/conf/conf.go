package conf

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

// EncodingKey : the encoding key.
var EncodingKey []byte

// IsDevMode : enable development mode.
var IsDevMode = false

// Create : create a config file.
func Create() error {
	data := map[string]any{
		"db_name":         "quid",
		"db_user":         "pguser",
		"db_password":     "my_password",
		"key":             generateRandomKey(),
		"enable_dev_mode": false,
	}

	jsonString, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile("config.json", jsonString, os.ModePerm)
}

// InitFromFile : get the config
// returns the postgres connection string.
func InitFromFile(isDevMode bool) (conn string, port int) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetDefault("db_name", "quid")
	viper.SetDefault("db_user", "pguser")
	viper.SetDefault("db_password", nil)
	viper.SetDefault("key", nil)
	viper.SetDefault("enable_dev_mode", false)

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("No config file found. Use the -conf option to generate one")
			os.Exit(4)
		}
		log.Fatal(err)
	}

	if isDevMode {
		IsDevMode = true
		if viper.Get("enable_dev_mode") == false {
			fmt.Println("Set enable_dev_mode to true in config in order to run in dev mode")
			os.Exit(1)
		}
	}

	hexKey := viper.Get("key").(string)
	if len(hexKey) < 32 {
		log.Panic("Want AES-128 key composed by 32 hexadecimal digits, but got ", len(hexKey))
	}
	EncodingKey, err = hex.DecodeString(hexKey[:32])
	if err != nil {
		fmt.Println("The key in config must be in hexadecimal format err=", err)
		os.Exit(5)
	}

	db := viper.Get("db_name").(string)
	usr := viper.Get("db_user").(string)
	pwd := viper.Get("db_password").(string)

	conn = "dbname=" + db + " user=" + usr + " password=" + pwd + " sslmode=disable"
	port = 8082
	return conn, port
}

func generateRandomKey() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}
