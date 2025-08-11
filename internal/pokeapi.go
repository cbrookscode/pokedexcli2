package internal

import (
	"encoding/json"
	"fmt"
	"io"
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
