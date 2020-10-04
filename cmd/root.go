/*Package cmd Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// devCmd represents the command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dev-tools",
	Short: "Dev tools for local environment",
	Long: `Dev Tools for local environment
		You can user one of the subcommands:
			- s3 (help to work with s3, e.g. backup, restore, list etc.)
			- dev (help to start services locally for development)
	`,
}

//

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/config.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
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
}
