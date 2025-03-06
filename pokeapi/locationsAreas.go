package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	cache "pokedex/pokecache"
)

type locationResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func DisplayLocationAreas(offset int, cache *cache.Cache) error {
	locationAreaUrl := fmt.Sprintf("https://pokeapi.co/api/v2/location-area?offset=%d&limit=20", offset)

	cachedData, found := cache.Get(locationAreaUrl)

	var locationAreasResponse locationResponse

	if found {
		fmt.Println("Using cached data...")

		decoder := json.NewDecoder(bytes.NewReader(cachedData))
		if err := decoder.Decode(&locationAreasResponse); err != nil {
			return fmt.Errorf("error decoding cached data: %v", err)
		}
	} else {
		fmt.Println("Fetching fresh data...")
		req, err := http.NewRequest("GET", locationAreaUrl, nil)
		if err != nil {
			return fmt.Errorf("%v", err)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("%v", err)
		}
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error reading response body: %v", err)
		}

		cache.Add(locationAreaUrl, bodyBytes)

		decoder := json.NewDecoder(bytes.NewReader(bodyBytes))
		if err := decoder.Decode(&locationAreasResponse); err != nil {
			return fmt.Errorf("error decoding response: %v", err)
		}
	}

	for _, location := range locationAreasResponse.Results {
		fmt.Println(location.Name)
	}

	return nil
}
