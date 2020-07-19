package config

import (
	"https://github.com/spf13/viper"
	"github.com/simplepki/core/keypair"
)

func init() {
	//from viper gh page
	viper.SetConfigName("settings")
	viper.AddConfigPath("/etc/simplepki/")
	viper.AddConfigPath("/opt/simplepki/")
	viper.AddConfigPath("$HOME/.simplepki/")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func IsCAEnabled() bool {
	return viper.IsSet("ca")
}

func GetCAStoreType() string {
	if viper.IsSet("ca.store") {
		switch viper.GetString("ca.store"){
		case "filesystem":
			return "filesystem"
		case "memory":
			return "memory"
		case "yubikey":
			return "yubikey"
		default:
			return "memory"
		}
	} else {
		return "filesystem"
	}
}

func ShouldOverwriteCA() bool {
	if viper.IsSet("ca.overwrite"){
		return viper.GetBool("ca.overwrite")
	}

	return false
}

func GetCAKeyPairConfig() (*keypair.KeyPairConfig, error) {
	config := &keypair.KeyPairConfig{}
	if viper.IsSet("ca.memory"){
		memConfig := GetInMemoryKeyPairConfig("ca.memory")

		config.KeyPairType = keypair.InMemory
		config.InMemoryConfig = memConfig
	} else if viper.IsSet("ca.filesystem"){
		fileConfig := &keypair.FileSystemKeyPairConfig{}

		config.KeyPairType = keypair.FileSyste,
		config.FileSystemConfig = fileConfig
	} else if viper.IsSet("ca.yubikey") {
		yubiConfig := &keypair.YubikeyKeyPairConfig{}

		config.KeyPairType = keypair.InMemory
		config.InMemoryConfig = memConfig
	} else {
		//default to memory
		memConfig := &keypair.InMemoryKeyPairConfig{}

		config.KeyPairType = keypair.InMemory
		config.InMemoryConfig = memConfig
	}


	return &keypair.KeyPairConfig{}, nil
}

func GetInMemoryKeyPairConfig (path string) *keypair.InMemoryKeyPairConfig {
	config := &keypair.InMemoryKeyPairConfig{}

	if viper.IsSet(path +".algorithm") {
		switch viper.IsSet(path +".algorithm"){
		case "ec256":
			config.KeyAgorithm = keypair.AlgorithmEC256
		case "ec384":
			config.KeyAgorithm = keypair.AlgorithmEC384
		case "rsa2048":
			config.KeyAgorithm = keypair.AlgorithmRSA2048
		case "rsa4096":
			config.KeyAgorithm = keypair.AlgorithmRSA4096
		}
	}
	return config
}