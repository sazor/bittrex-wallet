package client

import (
	"errors"
	"log"

	bittrex "github.com/sazor/go-bittrex"
	"github.com/theherk/viper"
)

func GetClient() (*bittrex.Bittrex, error) {
	key := viper.GetString("key")
	secret := viper.GetString("secret")
	if key == "" || secret == "" {
		log.Fatalln("Set API key and secret via config command.")
		return nil, errors.New("API config issue")
	}
	return bittrex.New(key, secret), nil
}
