package utils

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)


// IsExistCommand check if the program is in the PATH
func IsExistCommand(program string) bool {
	fmt.Printf("Check if %s command is exist\n", program)

	// Looking for the execution file with name in the PATH
	path, err := exec.LookPath(program)
	if err != nil {
		log.Fatalf("Installing %s is in your future\n", program)
		return false
	}
	fmt.Printf("%s is available at %s\n", program, path)
	return true
}

//ExecuteCommand execute program in the shell
func ExecuteCommand(command string, args ...string) error {

	cmd := exec.Command(command, args...)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	fmt.Printf("Get output of the command: %s\n", out.String())

	return err
}

// ExecuteCommandsWithPipe execute two commands with pipe command1 write output to the command2 as input
func ExecuteCommandsWithPipe(command1 string, args1 []string, command2 string, args2 []string) {
	c1 := exec.Command(command1, args1...)
	c2 := exec.Command(command2, args2...)

	r, w := io.Pipe()
	c1.Stdout = w
	c2.Stdin = r

	var b2 bytes.Buffer
	c2.Stdout = &b2

	c1.Start()
	c2.Start()
	c1.Wait()
	w.Close()
	c2.Wait()
	io.Copy(os.Stdout, &b2)
}
