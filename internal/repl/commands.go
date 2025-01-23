package repl

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/budgie-shell-go/internal/utils"
	"golang.org/x/sys/unix"
)

/*
	Internal commands for the shell.
	Things like exit, cd and type will be put here.
*/
var command_keys = []string {
	"echo", "exit", "pwd", "cd", "type",
}
var command_values = []func([]string) error {
	// ECHO
	func(args []string) error {
		if len(args) < 2 {
			fmt.Println("Usage: echo (message)")
			return errors.New("insufficient arguments")
		}

		message := strings.Join(args[1:], " ")
		fmt.Println(message)
		return nil
	},
	// EXIT
	func(args []string) error {
		if len(args) > 1 {
			exit_code_str := args[1]
			exit_code, err := strconv.Atoi(exit_code_str)
			if (err != nil) {
				fmt.Println("Provided invalid exit code: ", exit_code_str)
				return errors.New("invalid argument")
			}

			os.Exit(exit_code)
		} else {
			os.Exit(0)			
		}

		return nil
	},
	// PWD
	func(args []string) error {
		dir, err := os.Getwd()
		utils.IsNil(err)
		fmt.Println(dir)
		return nil
	},
	// CD
	func(args []string) error {
		if len(args) < 2 {
			fmt.Println("Usage: cd (directory)")
			return errors.New("insufficient arguments")
		}

		dir := args[1]
		if strings.HasPrefix(dir, "~") {
			dir = strings.Replace(dir, "~", home, 1)
		}	

		err := os.Chdir(dir)
		return err
	},
	// TYPE
	func(args []string) error {
		if len(args) < 2 {
			fmt.Println("Usage: type (command)")
			return errors.New("insufficient arguments")
		}

		cmd := args[1]
		for i := range(command_keys) {
			key := command_keys[i]
			if cmd == key {
				fmt.Println(cmd + " is a shell builtin")	
				return nil
			}
		}

		path, err := find_in_path(cmd)
		if err != nil {
			fmt.Println(cmd + ": not found")
			return errors.New("command not found")
		}

		fmt.Println(cmd + " is " + path)
		return nil
	},
}

/**
	 Find a specific file from the path variable.
	 If `command` contains slashes, it will attempt to find the file given absolute path.
*/
func find_in_path(command string) (string, error) {
	file, err := exec.LookPath(command)
	if err != nil {
		return "", err
	}

	fileInfo, err := os.Stat(file)
	if err != nil {
		return "", err
	}

	mode := fileInfo.Mode()
	if !(mode.IsRegular() || (uint32(mode & fs.ModeSymlink)) == 0) {
		return "", errors.New("File " + command + " is not a regular file or symlink.")
	}

	if (uint32(mode & 0111) == 0) {
		return "", errors.New("File " + command + " is not executable.")
	}

	if unix.Access(file, unix.X_OK) != nil {
		return "", errors.New("File " + command + " cannot be executed by this user.")
	}

	return file, nil
}

/**
	Execute a given command.
	`args[0]` will be treated as the command.
	Anything else provided in the args array will be treated as arguments to that command.
*/
func execute(args []string) error {
	var err error

	command := args[0]
	cmd := exec.Command(command, args[1:]...)
	
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()	
	if err != nil {
		fmt.Println("Recieved error: ", err)
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Println("Recieved error: ", err)
	}
	
	return nil
}
