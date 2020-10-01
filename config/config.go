package config

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/viper"
)

// Config reading the config file to continue
func Config(configFile string) {
	dPath := os.Getenv("DEV_CONFIG")
	// Build config path
	cPath := filepath.Join(dPath, "config.yaml")
	// Check if the file is exist
	if _, err := os.Stat(cPath); err == nil {
		fmt.Printf("File exists %s\n", cPath)
	} else {
		fmt.Printf("File does not exist %s\n", cPath)
	}

	// Get user home directory
	usr, error := user.Current()

	if error != nil {
		log.Fatal(error)
	}
	fmt.Println(usr.HomeDir)

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
	fmt.Println(viper.Get("dependency.bkms"))

	// Read properties file with aws settings
	viper.SetConfigType("properties")
	viper.SetConfigName("aws_credentials")

	// Merge all configs in one map
	err = viper.MergeInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: \n%s", err))
	}

	fmt.Println(viper.AllKeys())

	// Open the git repository
	r, err := git.PlainOpen(filepath.Join(usr.HomeDir, "workspace", "src/github.com/quiz"))

	// Read head of the current branch in the repository
	fmt.Println(r.Head())

	// Looking for the execution file with name in the PATH
	path, err := exec.LookPath("kubectl")
	if err != nil {
		log.Fatal("installing kubectl is in your future")
	}
	fmt.Printf("kubectl is available at %s\n", path)

	// Executing program in the shell
	cmd := exec.Command("ls")
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Get list of the files: %q\n", out.String())
}
