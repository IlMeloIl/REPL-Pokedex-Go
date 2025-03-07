package main

import (
	"bufio"
	"fmt"
	"os"
	api "pokedex/pokeapi"
	cache "pokedex/pokecache"
	"strings"
	"sync/atomic"
	"time"
)

type cliCommand struct {
	name        string
	description string
	callback    func([]string) error
}

var cliCommandsMap = map[string]cliCommand{}

func commandExit(s []string) error {
	fmt.Printf("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func commandHelp(s []string) error {
	fmt.Printf("Welcome to the Pokedex!\n")
	fmt.Printf("Usage:\n\n")
	for _, v := range cliCommandsMap {
		fmt.Printf("%v: %v\n", v.name, v.description)
	}
	return nil
}

var sharedOffset int32 = 0

func main() {
	cache := cache.NewCache(10 * time.Second)
	pokedex := make(map[string]api.Pokemon)
	cliCommandsMap = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays the names of the next 20 locations areas",
			callback: func(s []string) error {
				currentOffset := atomic.LoadInt32(&sharedOffset)
				err := api.DisplayLocationAreas(int(currentOffset), cache)
				if err == nil {
					atomic.AddInt32(&sharedOffset, 20)
				}
				return err
			},
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the names of the previous 20 locations areas",
			callback: func(s []string) error {
				currentOffset := atomic.LoadInt32(&sharedOffset)
				newOffset := currentOffset - 40
				if newOffset < 0 {
					atomic.StoreInt32(&sharedOffset, 0)
					fmt.Println("you're on the first page")
					return nil
				}
				err := api.DisplayLocationAreas(int(newOffset), cache)
				if err == nil {
					atomic.StoreInt32(&sharedOffset, newOffset)
				}
				return err
			},
		},
		"explore": {
			name:        "explore",
			description: "List all Pokemon in an area",
			callback: func(s []string) error {
				if len(s) < 2 {
					return fmt.Errorf("missing location area name. Usage: explore <area name>")
				}
				err := api.DisplayPokemonInArea(s[1], cache)
				return err
			},
		},
		"catch": {
			name:        "catch",
			description: "Try to catch a Pokemon",
			callback: func(s []string) error {
				if len(s) < 2 {
					return fmt.Errorf("missing Pokemon name. Usage: catch <Pokemon name>")
				}
				err := api.TryCatchPokemon(s[1], pokedex)
				return err
			},
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect Pokemon that's already caught",
			callback: func(s []string) error {
				if len(s) < 2 {
					return fmt.Errorf("missing Pokemon name. Usage: inspect <Pokemon name>")
				}
				pokemon, ok := pokedex[s[1]]
				if ok {
					fmt.Printf("Name: %s\n", pokemon.Name)
					fmt.Printf("Height: %d\n", pokemon.Height)
					fmt.Printf("Weight: %d\n", pokemon.Weight)
					fmt.Printf("Stats:\n")
					for k, v := range pokemon.Stats {
						fmt.Printf(" -%s: %d\n", k, v)
					}
					fmt.Printf("Types:\n")
					for _, type1 := range pokemon.Types {
						fmt.Printf(" - %s\n", type1)
					}
				} else {
					fmt.Println("you have not caught that Pokemon")
				}
				return nil
			},
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		err := scanner.Err()
		if err != nil {
			x := fmt.Errorf("%w", err)
			fmt.Println(x)
		}

		command := strings.Fields(strings.TrimSpace(scanner.Text()))

		if len(command) == 0 {
			continue
		}

		if cmd, ok := cliCommandsMap[command[0]]; ok {
			cmd.callback(command)

		} else {
			fmt.Printf("Unknown command\n")
		}
	}

}
