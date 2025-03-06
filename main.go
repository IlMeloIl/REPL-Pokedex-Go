package main

import (
	"bufio"
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var cliCommandsMap = map[string]cliCommand{}

func commandExit() error {
	fmt.Printf("Exiting the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Printf("Welcome to the Pokedex!\n")
	fmt.Printf("Usage:\n\n")
	for _, v := range cliCommandsMap {
		fmt.Printf("%v: %v\n", v.name, v.description)
	}
	return nil
}

func main() {

	cliCommandsMap = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}

	for {
		fmt.Print("Pokedex > ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		err := scanner.Err()
		if err != nil {
			x := fmt.Errorf("%w", err)
			fmt.Println(x)
		}

		command := scanner.Text()
		if cmd, ok := cliCommandsMap[command]; ok {
			cmd.callback()
		} else {
			fmt.Printf("Unknown command\n")
		}
	}

}
