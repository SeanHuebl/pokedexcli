package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/seanhuebl/pokedexcli/internal"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func getCommands(config *internal.Config, cache *internal.Cache, areaName, pokemonName string) map[string]cliCommand {

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
				if config == nil || cache == nil {
					return fmt.Errorf("config or cache not provided for 'map' command")
				}
				return internal.CommandMap(config, cache)
			},
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 location areas",
			callback: func() error {
				if config == nil || cache == nil {
					return fmt.Errorf("config or cache not provided for 'mapb' command")
				}
				return internal.CommandMapb(config, cache)
			},
		},
		"explore": {
			name:        "explore <location-area>",
			description: "Displays the pokemon in the area",
			callback: func() error {
				if cache == nil {
					return fmt.Errorf("cache not provided for 'explore' command")
				}
				return internal.Explore(cache, areaName)
			},
		},
		"catch": {
			name:        "catch <pokemon-name>",
			description: "Attempts to catch target pokemon",
			callback: func() error {
				if cache == nil {
					return fmt.Errorf("cache not provided for 'catch' command")
				}
				return internal.Catch(cache, pokemonName)
			},
		},
		"inspect": {
			name:        "inspect <pokemon-name>",
			description: "Displays information about caught pokemon",
			callback: func() error {
				return internal.Inspect(pokemonName)
			},
		},
		"pokedex": {
			name:        "pokedex",
			description: "display caught pokemon in pokedex",
			callback: func() error {
				return internal.ViewPokedex()
			},
		},
	}
}
func commandHelp() error {
	commands := getCommands(nil, nil, "", "")
	fmt.Println("Welcome to the Pokedex:")
	fmt.Println("Usage:")
	fmt.Println("")
	fmt.Printf("help: %v\n", commands["help"].description)
	fmt.Printf("map: %v\n", commands["map"].description)
	fmt.Printf("mapb: %v\n", commands["mapb"].description)
	fmt.Printf("explore <area-name>: %v\n", commands["explore"].description)
	fmt.Printf("catch <pokemon-name>: %v\n", commands["catch"].description)
	fmt.Printf("inspect <pokemon-name>: %v\n", commands["inspect"].description)
	fmt.Printf("pokedex: %v\n", commands["pokedex"].description)
	fmt.Printf("exit: %v\n\n", commands["exit"].description)
	return nil
}

func commandExit() error {
	os.Exit(0)
	return nil
}
func main() {
	const interval = 5 * time.Second
	ch := make(chan string)
	cache := internal.NewCache(interval)
	config := &internal.Conf

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
		helpMatch := regexp.MustCompile("help")
		exitMatch := regexp.MustCompile("exit")
		mapMatch := regexp.MustCompile("^map$")
		mapbMatch := regexp.MustCompile("mapb")
		exploreMatch := regexp.MustCompile("^explore .*")
		catchMatch := regexp.MustCompile("^catch .*")
		inspectMatch := regexp.MustCompile("^inspect .*")
		pokedexMatch := regexp.MustCompile("pokedex")
		switch true {

		case helpMatch.MatchString(value):
			commandHelp()
			fmt.Println("")

		case exitMatch.MatchString(value):
			commandExit()

		case mapMatch.MatchString(value):
			err := getCommands(config, cache, "", "")["map"].callback()
			if err != nil {
				fmt.Println(err)
			}
		case mapbMatch.MatchString(value):
			err := getCommands(config, cache, "", "")["mapb"].callback()
			if err != nil {
				fmt.Println(err)
			}
		case exploreMatch.MatchString(value):
			areaName := strings.Split(value, " ")[1]
			err := getCommands(nil, cache, areaName, "")["explore"].callback()
			if err != nil {
				fmt.Println(err)
			}
		case catchMatch.MatchString(value):
			pokemonName := strings.Split(value, " ")[1]
			err := getCommands(nil, cache, "", pokemonName)["catch"].callback()
			if err != nil {
				fmt.Println(err)
			}
		case inspectMatch.MatchString(value):
			pokemonName := strings.Split(value, " ")[1]
			err := getCommands(nil, cache, "", pokemonName)["inspect"].callback()
			if err != nil {
				fmt.Println(err)
			}
		case pokedexMatch.MatchString(value):
			err := getCommands(nil, nil, "", "")["pokedex"].callback()
			if err != nil {
				fmt.Println(err)
			}

		}

	}
}
