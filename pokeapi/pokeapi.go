package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	cache "pokedex/pokecache"
	"time"
)

type locationResponseStruct struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type pokemonEncountersStruct struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	Name   string
	Height int
	Weight int
	Stats  map[string]int
	Types  []string
}

type catchPokemonStruct struct {
	Height int `json:"height"`
	Weight int `json:"weight"`
	Stats  []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
	BaseExperience int `json:"base_experience"`
}

func fetchFromApi(url string, target interface{}, cache *cache.Cache, useCache bool) error {
	if useCache && cache != nil {
		cachedData, found := cache.Get(url)
		if found {
			return json.NewDecoder(bytes.NewReader(cachedData)).Decode(target)
		}
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found")
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	if useCache && cache != nil {
		cache.Add(url, bodyBytes)
	}

	return json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(target)
}

func DisplayLocationAreas(offset int, cache *cache.Cache) error {
	show20AreasUrl := fmt.Sprintf("https://pokeapi.co/api/v2/location-area?offset=%d&limit=20", offset)

	var locationAreasResponse locationResponseStruct

	err := fetchFromApi(show20AreasUrl, &locationAreasResponse, cache, true)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	for _, location := range locationAreasResponse.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func DisplayPokemonInArea(locationAreaName string, cache *cache.Cache) error {
	locationAreaUrl := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", locationAreaName)

	var pokemonEncounters pokemonEncountersStruct

	err := fetchFromApi(locationAreaUrl, &pokemonEncounters, cache, true)
	if err != nil {
		fmt.Println("Area not found")
		return fmt.Errorf("%w", err)
	}

	fmt.Printf("Exploring %s...\n", locationAreaName)
	for _, pokemon := range pokemonEncounters.PokemonEncounters {
		fmt.Printf(" - %s\n", pokemon.Pokemon.Name)
	}

	return nil
}

func attemptCatch(baseExperience int) bool {
	catchChange := 100.0 - (float64(baseExperience) * 0.4)

	if catchChange < 5.0 {
		catchChange = 5.0
	}
	if catchChange > 90.0 {
		catchChange = 90.0
	}

	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	roll := rng.Float64() * 100

	return roll <= catchChange
}

func TryCatchPokemon(pokemonName string, pokedex map[string]Pokemon) error {
	pokemonUrl := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", pokemonName)

	var pokemon catchPokemonStruct

	err := fetchFromApi(pokemonUrl, &pokemon, nil, false)
	if err != nil {
		fmt.Println("Pokemon not found")
		return fmt.Errorf("Pokemon not found")
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)
	if attemptCatch(pokemon.BaseExperience) {
		fmt.Printf("%s was caught!\n", pokemonName)
		statsMap := make(map[string]int)
		for _, stat := range pokemon.Stats {
			statsMap[stat.Stat.Name] = stat.BaseStat
		}

		types := make([]string, 0, len(pokemon.Types))
		for _, t := range pokemon.Types {
			types = append(types, t.Type.Name)
		}

		pokedex[pokemonName] = Pokemon{Name: pokemonName, Height: pokemon.Height, Weight: pokemon.Weight, Stats: statsMap, Types: types}
		fmt.Println("You may now inspect it with the inspect command.")
	} else {
		fmt.Printf("%s escaped!\n", pokemonName)
	}

	return nil

}
