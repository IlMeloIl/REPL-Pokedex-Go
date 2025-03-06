package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	for {
		fmt.Print("Pokedex > ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		err := scanner.Err()
		if err != nil {
			x := fmt.Errorf("%w", err)
			fmt.Println(x)
		}

		cleanedString := strings.Fields(strings.ToLower(scanner.Text()))
		fmt.Println(cleanedString[0])
	}

}
