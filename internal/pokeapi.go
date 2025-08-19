package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LocationNames struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type ListofLocations struct {
	Count    int             `json:"count"`
	Next     string          `json:"next"`
	Previous any             `json:"previous"`
	Results  []LocationNames `json:"results"`
}

type LocationArea struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int           `json:"chance"`
				ConditionValues []interface{} `json:"condition_values"`
				MaxLevel        int           `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func GetLocations(url string) (ListofLocations, []byte, error) {
	areas := ListofLocations{}

	resp, err := http.Get(url)
	if err != nil {
		return ListofLocations{}, nil, fmt.Errorf("error with GET request at %v", url)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return ListofLocations{}, nil, fmt.Errorf("error reading http get reponse: %v", err)
	}
	err = json.Unmarshal(data, &areas)
	if err != nil {
		return ListofLocations{}, nil, fmt.Errorf("error unmarshalling reponse from GET request at %v", url)
	}

	return areas, data, nil
}

func GetPokemon(location string) (LocationArea, error) {
	url := "https://pokeapi.co/api/v2/location-area/" + location
	area := LocationArea{}

	resp, err := http.Get(url)
	if err != nil {
		return area, fmt.Errorf("error with GET request: %v", err)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return area, fmt.Errorf("error reading http get response: %v", err)
	}

	err = json.Unmarshal(data, &area)
	if err != nil {
		return area, fmt.Errorf("error unmarshalling json: %v", err)
	}

	return area, nil
}
