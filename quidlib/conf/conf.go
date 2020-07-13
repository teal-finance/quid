package conf

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

var dbName string
var dbUser string
var dbPassword string

// ConnStr : the postgres connection string
var ConnStr string

// EncodingKey : the encoding key
var EncodingKey string

// DefaultTokenTimeout : default timeout for user tokens
var DefaultTokenTimeout time.Duration

// Create : create a config file
func Create() {
	data := map[string]interface{}{
		"db_name":                "quid",
		"db_user":                "pguser",
		"db_password":            "",
		"default_tokens_timeout": "24h",
		"key":                    genKey(),
	}
	jsonString, _ := json.MarshalIndent(data, "", "    ")
	ioutil.WriteFile("config.json", jsonString, os.ModePerm)
}

// Init : get the config
func Init() (bool, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetDefault("db_name", "quid")
	viper.SetDefault("db_user", "pguser")
	viper.SetDefault("db_password", nil)
	viper.SetDefault("default_tokens_timeout", "24h")
	viper.SetDefault("key", nil)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found
			return false, nil
		}
		log.Fatal(err)
	}
	dbName = viper.Get("db_name").(string)
	dbUser = viper.Get("db_user").(string)
	dbPassword = viper.Get("db_password").(string)
	var err error
	DefaultTokenTimeout, err = time.ParseDuration(viper.Get("default_tokens_timeout").(string))
	if err != nil {
		return true, err
	}
	EncodingKey = viper.Get("key").(string)
	ConnStr = "user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " sslmode=disable"
	return true, nil
}

// GenKey : generate a random key
func genKey() string {
	b, err := generateRandomBytes(32)
	if err != nil {
		log.Fatal(err)
	}
	h := hmac.New(sha256.New, []byte(b))
	token := hex.EncodeToString(h.Sum(nil))
	return token

}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
