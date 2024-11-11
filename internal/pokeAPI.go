package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

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
