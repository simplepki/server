package config

import (
	"github.com/spf13/viper"
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

func GetFileSystemKeyPairConfig (path string) *keypair.FileSystemKeyPairConfig {
	config := &keypair.FileSystemKeyPairConfig{}

	/*if viper.IsSet(path +".algorithm") {
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
	}*/

	if viper.IsSet(path + ".key_file") {
		config.KeyFile = viper.GetString(path + ".key_file")
	} else {
		config.KeyFile = "./key.pem"
	}

	if viper.IsSet(path + ".cert_file") {
		config.KeyFile = viper.GetString(path + ".cert_file")
	} else {
		config.KeyFile = "./cert.pem"
	}

	if viper.IsSet(path + ".chain_file") {
		config.KeyFile = viper.GetString(path + ".chain_file")
	} else {
		config.KeyFile = "./chain.pem"
	}
	return config
}

func GetYubikeyKeyPairConfig (path string) *keypair.YubikeyKeyPairConfig {
	config := &keypair.YubikeyKeyPairConfig{}

	if viper.IsSet(path + ".subject_name") {
		config.CertSubjectName = viper.GetString(path + ".subject_name")
	} 

	if viper.IsSet(path + ".reset") {
		config.Reset = viper.GetBool(path + ".reset")
	}

	if viper.IsSet(path+".yubikey_name") {
		name := viper.GetString("")
		config.Name = &name
	}

	if viper.IsSet(path+".yubikey_serial_number") {
		num := viper.GetUint32(path+".yubikey_serial_number")
		config.Serial = &num
	}

	if viper.IsSet(path+".pin") {
		pin := viper.GetString(path+".pin")
		config.PIN = &pin
	}

	if viper.IsSet(path+".puk") {
		puk := viper.GetString(path+".puk")
		config.PUK = &puk
	}

	if viper.IsSet(path+".management_key") {
		mk := viper.GetString(path+".management_key")
		config.Base64ManagementKey = &mk
	}
 	return config
}