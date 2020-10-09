package config

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config return a map with configuration properties and secrets
func Config() *viper.Viper {
	// Get user home directory
	usr, error := user.Current()

	if error != nil {
		log.Fatal(error)
	}
	fmt.Println(usr.HomeDir)

	dPath, exist := os.LookupEnv("DEV_CONFIG")

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
	viper.AddConfigPath(filepath.Join(usr.HomeDir, ".aws"))

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: \n%s", err))
	}

	// Read properties file with aws settings
	viper.SetConfigType("properties")
	viper.SetConfigName("credentials")

	// Merge all configs in one map
	err = viper.MergeInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: \n%s", err))
	}
	return viper.GetViper()
}
