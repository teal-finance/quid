package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"os"

	"github.com/spf13/viper"
)

// createConfigFile : createConfigFile a config file.
func createConfigFile(dbName, dbUser, dbPass string) error {
	data := map[string]any{
		"db_name":     dbName,
		"db_user":     dbUser,
		"db_password": dbPass,
		"key":         randomAES128KeyHex(),
	}

	jsonString, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile("config.json", jsonString, os.ModePerm)
}

// readConfigFile : get the config.
func readConfigFile() (name, usr, pwd, key string) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetDefault("db_name", nil)
	viper.SetDefault("db_user", nil)
	viper.SetDefault("db_password", nil)
	viper.SetDefault("key", nil)

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Info(`No "config.json" file found. You may use the -conf flag to generate the "config.json" file with a random key.`)
		} else {
			log.Fatal(err)
		}
	}

	name = viper.Get("db_name").(string)
	usr = viper.Get("db_user").(string)
	pwd = viper.Get("db_password").(string)
	key = viper.Get("key").(string)

	// conn = "dbname=" + name + " user=" + usr + " password=" + pwd + " sslmode=disable"
	return name, usr, pwd, key
}

func randomAES128KeyHex() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		log.Panic(err)
	}
	return hex.EncodeToString(bytes)
}
