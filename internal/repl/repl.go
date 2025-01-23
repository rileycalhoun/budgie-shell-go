package repl

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/budgie-shell/internal/utils"
)

var home = os.Getenv("HOME")

func Read() []string {
	// Create reader to from stdin
	reader := bufio.NewReader(os.Stdin)

	// Get current working directory
	cwd, err := os.Getwd()

	// Ensure there is no error; if there is then panic
	utils.IsNil(err)

	// Check if the cwd starts with the user's home directory
	if strings.HasPrefix(cwd, home) {
		// If it does, replace it with a ~
		cwd = strings.Replace(cwd, home, "~", 1)
	}

	// Print the cwd and a $ afterwards
	fmt.Printf("%s$ ", cwd)

	// Read string, end when a newline is formed
	text, err := reader.ReadString('\n')
	utils.IsNil(err)	

	// Trim command
	text = strings.Trim(text, " \t\n")

	// Split by spaces; eventually allow for quotes
	args := strings.Split(text, " ");

	return args
}

func Eval(args []string) error {
	command := args[0]

	for i := range(command_keys) {
		name := command_keys[i]
		if name == command {
			return command_values[i](args)
		}
	}

	_, err := find_in_path(command)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return execute(args)
}

