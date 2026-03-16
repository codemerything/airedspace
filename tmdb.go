package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type TMDB struct {
	client  *http.Client
	apiKey  string
	baseURL string
}

type TMDBResult struct {
	ID int `json:"id"`
}

type TMDBResponse struct {
	Results []TMDBResult `json:"results"`
}

func NewTMDB(apiKey, baseURL string) *TMDB {
	return &TMDB{
		client:  &http.Client{},
		apiKey:  apiKey,
		baseURL: baseURL,
	}

}

// fetchmoviebyid gets the movieid by searching for the movie by name which
// will extract the movies id

func (t *TMDB) FetchMovieID(search string) (int, error) {
	req, err := http.NewRequest("GET", t.baseURL+"search/movie?query="+search, nil)
	if err != nil {
		return 0, fmt.Errorf("Failed request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer"+t.apiKey)

	resp, err := t.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("search failed: %w", err)
	}

	defer resp.Body.Close()

	var data TMDBResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(data.Results) == 0 {
		return 0, errors.New("no films found")
	}

	return data.Results[0].ID, nil
}

func (t *TMDB) FetchDetails(id int) (Films, error) {
	// start with the request
	req, err := http.NewRequest("GET", t.baseURL+"movie/"+strconv.Itoa(id), nil)

	if err != nil {
		return Films{}, errors.New("Request failed")
	}

	req.Header.Set("Authorization", "Bearer"+t.apiKey)
	resp, err := t.client.Do(req)
	if err != nil {
		return Films{}, fmt.Errorf("search failed: %w", err)
	}
}
