package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
)

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
	CatchChance            float32
}
type Ability struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Abilities struct {
	Ability  Ability `json:"ability"`
	IsHidden bool    `json:"is_hidden"`
	Slot     int     `json:"slot"`
}
type Move struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Moves struct {
	Move Move `json:"move"`
}
type Species struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Stat struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Stats struct {
	BaseStat int  `json:"base_stat"`
	Effort   int  `json:"effort"`
	Stat     Stat `json:"stat"`
}
type Type struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Types struct {
	Slot int  `json:"slot"`
	Type Type `json:"type"`
}

type Location struct {
	Count    int     `json:"count"`
	Next     string  `json:"next"`
	Previous *string `json:"previous"`
	Results  []Area  `json:"results"`
}
type Area struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type PokemonEncounters []struct {
	Pokemon struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"pokemon"`
}

type Config struct {
	Next     string
	Previous *string
}

var Conf = Config{
	Next:     "https://pokeapi.co/api/v2/location-area/",
	Previous: nil,
}

var Pokedex = make(map[string]Pokemon)

func CommandMap(cF *Config, c *Cache) error {
	var location Location
	v, ok := c.entry[cF.Next]
	if !ok {

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
		c.Add(cF.Next, data)

		if err := json.Unmarshal(data, &location); err != nil {
			return err
		}
	} else {
		if err := json.Unmarshal(v.val, &location); err != nil {
			return err
		}
	}

	cF.Next = location.Next
	cF.Previous = location.Previous

	for _, result := range location.Results {
		fmt.Println(result.Name)
	}
	return nil
}

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

	for _, result := range location.Results {
		fmt.Println(result.Name)
	}

	return nil
}

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

	if pokemon.BaseExperience < 70 {
		pokemon.CatchChance = .85
	} else if pokemon.BaseExperience < 160 {
		pokemon.CatchChance = .50
	} else if pokemon.BaseExperience < 250 {
		pokemon.CatchChance = .25
	} else {
		pokemon.CatchChance = .10
	}

	catchRoll := rand.Float32()
	if catchRoll < pokemon.CatchChance {
		fmt.Printf("%s was caught!\n\n", pokemonName)
		Pokedex[pokemonName] = pokemon
		return nil
	}
	fmt.Printf("%s escaped!\n\n", pokemonName)
	return nil
}
func Inspect(pokemonName string) error {
	val, ok := Pokedex[pokemonName]
	if !ok {
		return fmt.Errorf("pokemon not found")
	}
	fmt.Printf("Name: %v\nHeight: %v\nWeight: %v\nStats:\n", val.Name, val.Height, val.Weight)
	for _, stat := range val.Stats {
		fmt.Printf("  -%v: %v\n", stat.Stat.Name, stat.BaseStat)
	}
	return nil
}
