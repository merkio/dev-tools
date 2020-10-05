package utils

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

func ListPods() {
	fmt.Println("List pods for env")
}

func ListDB() {
	fmt.Println("List databases for env")
}

func IsExistCommand(program string) {
	fmt.Printf("Check if %s command is exist\n", program)

	// Looking for the execution file with name in the PATH
	path, err := exec.LookPath(program)
	if err != nil {
		log.Fatal("Installing %s is in your future\n", program)
	}
	fmt.Printf("%s is available at %s\n", program, path)

}

func RunService(name string, profile string, mode string) {
	fmt.Printf("Run service %s in mode %s with profile %s\n", name, mode, profile)
}

func CreateDatabase(service string, dbName string) {

}

func executeCommand(command string, args ...string) {

	// Executing program in the shell
	cmd := exec.Command(command, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Get output of the command: %s\n", out.String())

}
