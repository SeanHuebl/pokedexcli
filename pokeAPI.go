package main

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
	Results  []Areas `json:"results"`
}
type Areas struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Config struct {
	Next     string
	Previous *string
}

var config = Config{
	Next:     "https://pokeapi.co/api/v2/location-area/",
	Previous: nil,
}

func commandMap(cF *Config) error {
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
	var location Location
	if err = json.Unmarshal(data, &location); err != nil {
		return err
	}

	cF.Next = location.Next
	cF.Previous = location.Previous

	for _, result := range location.Results {
		fmt.Println(result.Name)
	}
	return nil
}

func commandMapb(cF *Config) error {
	if cF.Previous == nil {
		return fmt.Errorf("no previous location")
	}
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
	var location Location
	if err = json.Unmarshal(data, &location); err != nil {
		return err
	}

	cF.Next = location.Next
	cF.Previous = location.Previous

	for _, result := range location.Results {
		fmt.Println(result.Name)
	}

	return nil
}
