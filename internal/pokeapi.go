package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type LocationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type ListofLocations struct {
	Count    int            `json:"count"`
	Next     string         `json:"next"`
	Previous any            `json:"previous"`
	Results  []LocationArea `json:"results"`
}

func GetLocations(url string) (ListofLocations, error) {
	areas := ListofLocations{}

	resp, err := http.Get(url)
	if err != nil {
		return ListofLocations{}, fmt.Errorf("error with GET request at %v", url)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&areas)
	if err != nil {
		return ListofLocations{}, fmt.Errorf("error unmarshalling reponse from GET request at %v", url)
	}

	return areas, nil
}
