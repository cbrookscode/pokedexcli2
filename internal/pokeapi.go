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

type Pokemon struct {
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
	} `json:"abilities"`
	BaseExperience         int    `json:"base_experience"`
	Height                 int    `json:"height"`
	ID                     int    `json:"id"`
	IsDefault              bool   `json:"is_default"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Moves                  []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt  int `json:"level_learned_at"`
			MoveLearnMethod struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"move_learn_method"`
			Order        interface{} `json:"order"`
			VersionGroup struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version_group"`
		} `json:"version_group_details"`
	} `json:"moves"`
	Name    string `json:"name"`
	Order   int    `json:"order"`
	Species struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
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

func PrintLocations(url string) (ListofLocations, error) {
	locations, _, err := GetLocations(url)
	if err != nil {
		return locations, fmt.Errorf("error making new get request: %v", err)
	}

	fmt.Println("-------------------------------")
	for _, area := range locations.Results {
		fmt.Println(area.Name)
	}
	fmt.Println("-------------------------------")

	return locations, nil
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

func CalcChancetoCatch(baseXP int) float64 {
	// pick number between 1 - 100 this is number to meet or beat
	// total xp granted is per pokemon variant of difficulty to catch
	/*
		scale determines effectives of xp as difficulty scaler. lower number means xp has higher impact
		base chance = 100/ (100/(1+ (xp/scale)))
		Convet base chance to odds with below
		odds = bc/(1 - bc)
		multiply your other multiplers that would increase chance against odds
		go back to proability
		p = odds / (1 + odds)
		multiply p by 100 to get the integer for comparison for the random number between 1 - 100
	*/
}
