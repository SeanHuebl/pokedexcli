package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
)

// Pokemon represents a Pokémon entity with various attributes such as abilities, stats, and encounter information.
type Pokemon struct {
	Abilities              []Abilities `json:"abilities"`
	BaseExperience         int         `json:"base_experience"`
	Height                 int         `json:"height"`
	HeldItems              []any       `json:"held_items"`
	ID                     int         `json:"id"`
	LocationAreaEncounters string      `json:"location_area_encounters"`
	Moves                  []Moves     `json:"moves"`
	Name                   string      `json:"name"`
	Species                Species     `json:"species"`
	Stats                  []Stats     `json:"stats"`
	Types                  []Types     `json:"types"`
	Weight                 int         `json:"weight"`
	CatchChance            float32     // Calculated chance of catching the Pokémon.
}

// Ability represents a specific ability of a Pokémon.
type Ability struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Abilities holds an ability, its slot, and a hidden status flag for a Pokémon.
type Abilities struct {
	Ability  Ability `json:"ability"`
	IsHidden bool    `json:"is_hidden"`
	Slot     int     `json:"slot"`
}

// Move represents a specific move of a Pokémon.
type Move struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Moves contains the information of a Pokémon's move.
type Moves struct {
	Move Move `json:"move"`
}

// Species represents the species information of a Pokémon.
type Species struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Stat provides details about a particular stat of a Pokémon.
type Stat struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Stats contains the stat details for a Pokémon.
type Stats struct {
	BaseStat int  `json:"base_stat"`
	Effort   int  `json:"effort"`
	Stat     Stat `json:"stat"`
}

// Type represents a type attribute for a Pokémon.
type Type struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Types holds type-related data for a Pokémon.
type Types struct {
	Slot int  `json:"slot"`
	Type Type `json:"type"`
}

// Location represents location data with count, pagination, and location results.
type Location struct {
	Count    int     `json:"count"`
	Next     string  `json:"next"`
	Previous *string `json:"previous"`
	Results  []Area  `json:"results"`
}

// Area provides details about a specific location area.
type Area struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

// PokemonEncounters holds a list of Pokémon that can be encountered in a specific area.
type PokemonEncounters []struct {
	Pokemon struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"pokemon"`
}

// Config contains configuration options for navigating between Pokémon locations.
type Config struct {
	Next     string
	Previous *string
}

// Conf initializes the configuration for location-based API navigation.
var Conf = Config{
	Next:     "https://pokeapi.co/api/v2/location-area/",
	Previous: nil,
}

// Pokedex is an in-memory storage for caught Pokémon data.
var Pokedex = make(map[string]Pokemon)

// CommandMap retrieves and processes the next set of location results, storing them in the cache.
func CommandMap(cF *Config, c *Cache) error {
	var location Location
	v, ok := c.entry[cF.Next]
	if !ok {
		// Fetch data from API if not in cache.
		client := &http.Client{}
		req, err := http.NewRequest("GET", cF.Next, nil)
		if err != nil {
			log.Fatal(err)
		}

		res, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		c.Add(cF.Next, data) // Add data to cache

		if err := json.Unmarshal(data, &location); err != nil {
			return err
		}
	} else {
		// Use cached data if available.
		if err := json.Unmarshal(v.val, &location); err != nil {
			return err
		}
	}

	cF.Next = location.Next
	cF.Previous = location.Previous

	// Display each location name.
	for _, result := range location.Results {
		fmt.Println(result.Name)
	}
	return nil
}

// CommandMapb retrieves the previous set of location results, using cached data if available.
func CommandMapb(cF *Config, c *Cache) error {
	if cF.Previous == nil {
		return fmt.Errorf("no previous location")
	}
	var location Location
	v, ok := c.entry[*cF.Previous]
	if !ok {
		client := &http.Client{}
		req, err := http.NewRequest("GET", *cF.Previous, nil)
		if err != nil {
			log.Fatal(err)
		}

		res, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		c.Add(*cF.Previous, data)

		if err = json.Unmarshal(data, &location); err != nil {
			return err
		}
	} else {
		if err := json.Unmarshal(v.val, &location); err != nil {
			return err
		}
	}

	cF.Next = location.Next
	cF.Previous = location.Previous

	// Display each location name.
	for _, result := range location.Results {
		fmt.Println(result.Name)
	}

	return nil
}

// Explore retrieves Pokémon encounter data for a specified area and displays encountered Pokémon.
func Explore(c *Cache, areaName string) error {
	fmt.Printf("Exploring %s...\nPokemon:\n", areaName)
	areaUrl := "https://pokeapi.co/api/v2/location-area/" + areaName
	var encounters struct {
		PokemonEncounters PokemonEncounters `json:"pokemon_encounters"`
	}
	v, ok := c.entry[areaName]
	if !ok {
		client := &http.Client{}
		req, err := http.NewRequest("GET", areaUrl, nil)
		if err != nil {
			log.Fatal(err)
		}

		res, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		c.Add(areaUrl, data)

		if err = json.Unmarshal(data, &encounters); err != nil {
			return err
		}
	} else {
		if err := json.Unmarshal(v.val, &encounters); err != nil {
			return err
		}
	}

	for _, encounter := range encounters.PokemonEncounters {
		fmt.Printf("- %s\n", encounter.Pokemon.Name)
	}
	return nil
}

// Catch attempts to catch a specified Pokémon, calculating success based on its experience level.
func Catch(c *Cache, pokemonName string) error {
	pokemonUrl := "https://pokeapi.co/api/v2/pokemon/" + pokemonName
	fmt.Printf("Throwing a ball at %s...\n", pokemonName)
	var pokemon Pokemon
	v, ok := c.entry[pokemonName]
	if !ok {
		client := &http.Client{}
		req, err := http.NewRequest("GET", pokemonUrl, nil)
		if err != nil {
			log.Fatal(err)
		}

		res, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		c.Add(pokemonUrl, data)

		if err = json.Unmarshal(data, &pokemon); err != nil {
			return err
		}
	} else {
		if err := json.Unmarshal(v.val, &pokemon); err != nil {
			return err
		}
	}

	// Calculate the catch chance based on Pokémon experience.
	if pokemon.BaseExperience < 70 {
		pokemon.CatchChance = .85
	} else if pokemon.BaseExperience < 160 {
		pokemon.CatchChance = .50
	} else if pokemon.BaseExperience < 250 {
		pokemon.CatchChance = .25
	} else {
		pokemon.CatchChance = .10
	}

	// Attempt to catch the Pokémon.
	catchRoll := rand.Float32()
	if catchRoll < pokemon.CatchChance {
		fmt.Printf("%s was caught!\ndata added to pokedex\nto view data use the inspect command\n\n", pokemonName)
		Pokedex[pokemonName] = pokemon
		return nil
	}
	fmt.Printf("%s escaped!\n\n", pokemonName)
	return nil
}

// Inspect displays detailed information about a specified Pokémon in the Pokedex.
func Inspect(pokemonName string) error {
	val, ok := Pokedex[pokemonName]
	if !ok {
		return fmt.Errorf("pokemon not found\n")
	}
	fmt.Printf("Name: %v\nHeight: %v\nWeight: %v\nStats:\n", val.Name, val.Height, val.Weight)
	for _, stat := range val.Stats {
		fmt.Printf("  - %v: %v\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, poketype := range val.Types {
		fmt.Printf("  - %v\n", poketype.Type.Name)
	}
	return nil
}

// ViewPokedex displays all Pokémon currently in the Pokedex.
func ViewPokedex() error {
	fmt.Println("Your pokedex:")
	if len(Pokedex) == 0 {
		fmt.Printf("Your pokedex contains 0 pokemon\n\n")
	} else {
		for key := range Pokedex {
			fmt.Printf("  - %v\n", key)
		}
	}
	fmt.Println("")
	return nil
}
