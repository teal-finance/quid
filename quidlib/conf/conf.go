package conf

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/viper"
)

// ConnStr : the postgres connection string.
var ConnStr string

// EncodingKey : the encoding key.
var EncodingKey string

// IsDevMode : enable development mode.
var IsDevMode = false

// Port : the port to run on.
var Port = "8082"

// Create : create a config file.
func Create() error {
	data := map[string]interface{}{
		"db_name":         "quid",
		"db_user":         "pguser",
		"db_password":     "my_password",
		"key":             generateRandomKey(),
		"enable_dev_mode": false,
	}

	jsonString, _ := json.MarshalIndent(data, "", "    ")

	return ioutil.WriteFile("config.json", jsonString, os.ModePerm)
}

// InitFromFile : get the config.
func InitFromFile(isDevMode bool) (bool, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetDefault("db_name", "quid")
	viper.SetDefault("db_user", "pguser")
	viper.SetDefault("db_password", nil)
	viper.SetDefault("key", nil)
	viper.SetDefault("enable_dev_mode", false)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found
			return false, nil
		}
		log.Fatal(err)
	}

	if isDevMode {
		if viper.Get("enable_dev_mode") == false {
			fmt.Println("Please set enable_dev_mode to true in config in order to run in dev mode")
			os.Exit(1)
		}
	}

	dbname := viper.Get("db_name").(string)
	user := viper.Get("db_user").(string)
	password := viper.Get("db_password").(string)
	ConnStr = "dbname=" + dbname + " user=" + user + " password=" + password + " sslmode=disable"

	EncodingKey = viper.Get("key").(string)

	return true, nil
}

func generateRandomKey() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		panic(err.Error())
	}

	return hex.EncodeToString(bytes)
}
