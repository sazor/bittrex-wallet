package config

import (
	"log"

	"github.com/mitchellh/go-homedir"
	"github.com/theherk/viper"
)

func Load(cfgFile string) error {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		err := newCfgFile()
		if err != nil {
			return err
		}
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln("Cant read config")
		return err
	}
	log.Println("Using config file:", viper.ConfigFileUsed())
	return nil
}

func newCfgFile() error {
	// Find home directory
	home, err := homedir.Dir()
	if err != nil {
		log.Fatalln(err)
		return err
	}

	// Search config in home directory with name ".bittrex-wallet" (without extension)
	viper.AddConfigPath(home)
	viper.SetConfigName(".bittrex-wallet")
	return nil
}
