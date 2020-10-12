package config

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/spf13/viper"
)

// ConfigMap presentation of config file as a map with maps
type ConfigMap struct {
	Repositories map[string]map[string]string `mapstructure:"repositories"`
	Dependency   map[string][]string          `mapstructure:"dependency"`
}

// Config return a map with configuration properties and secrets
func Config() *viper.Viper {
	// Get user home directory
	usr, error := user.Current()

	if error != nil {
		log.Fatal(error)
	}
	fmt.Println(usr.HomeDir)

	dPath, exist := os.LookupEnv("DEV_CONFIG")

	viper.Set("BackupDir", filepath.Join(usr.HomeDir, "bkms-backup"))

	if !exist {
		dPath = usr.HomeDir
	}
	// Build config path
	cPath := filepath.Join(dPath, "config.yaml")
	// Check if the file is exist
	if _, err := os.Stat(cPath); err == nil {
		fmt.Printf("File exists %s\n", cPath)
	} else {
		fmt.Printf("File does not exist %s\n", cPath)
	}

	// Read configs
	// Read yaml
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(dPath)
	viper.AddConfigPath(usr.HomeDir)

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: \n%s", err))
	}

	awsConf := filepath.Join(usr.HomeDir, ".aws")
	// Check if the file is exist
	if _, err := os.Stat(awsConf); err == nil {
		fmt.Printf("File exists %s\n", awsConf)

		// Read properties file with aws settings
		viper.AddConfigPath(awsConf)
		viper.SetConfigType("properties")
		viper.SetConfigName("credentials")

		// Merge all configs in one map
		err = viper.MergeInConfig() // Find and read the config file
		if err != nil {             // Handle errors reading the config file
			panic(fmt.Errorf("Fatal error config file: \n%s", err))
		}
	}
	return viper.GetViper()
}

// GetConfigMap get configuration map
func GetConfigMap() ConfigMap {
	Config()
	configMap := &ConfigMap{}
	err := viper.Unmarshal(configMap)

	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	return *configMap
}
