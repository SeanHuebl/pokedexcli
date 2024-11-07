package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/seanhuebl/pokedexcli/internal"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func getCommands() map[string]cliCommand {

	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},

		"exit": {
			name:        "exit",
			description: "Exits the program",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Displays the next 20 location areas",
			callback: func() error {
				return internal.CommandMap(&internal.Conf)
			},
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 location areas",
			callback: func() error {
				return internal.CommandMapb(&internal.Conf)
			},
		},
	}
}

func commandHelp() error {
	commands := getCommands()
	fmt.Println("Welcome to the Pokedex:")
	fmt.Println("Usage:")
	fmt.Println("")
	fmt.Printf("help: %v\n", commands["help"].description)
	fmt.Printf("exit: %v\n", commands["exit"].description)
	return nil
}

func commandExit() error {
	os.Exit(0)
	return nil
}
func main() {

	ch := make(chan string)
	go func() {
		defer close(ch)
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			ch <- scanner.Text()
		}
	}()

	for {
		fmt.Print("pokedex > ")
		value := <-ch
		fmt.Println("")
		switch value {

		case "help":
			commandHelp()
			fmt.Println("")

		case "exit":
			commandExit()

		case "map":
			err := getCommands()["map"].callback()
			if err != nil {
				fmt.Println(err)
			}
		case "mapb":
			err := getCommands()["mapb"].callback()
			if err != nil {
				fmt.Println(err)
			}
			
		}

	}
}
