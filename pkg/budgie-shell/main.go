/*
The Budgie Shell is a very simple implementation of a POSIX-compliant shell.
I created this to help myself learn Go, and to deepen my understanding of Unix.
*/
package main

import "github.com/budgie-shell/internal/repl"

var logs []string = make([]string, 0) 

func main() {
	// Read any configuration files; set aliases, etc.

	for {
		args := repl.Read();
		err := repl.Eval(args)
		if err != nil {
			logs = append(logs, err.Error())
		}
	}
	
}
