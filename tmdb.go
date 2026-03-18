package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type TMDB struct {
	client  *http.Client
	apiKey  string
	baseURL string
}

type TMDBFilmDetails struct {
	FilmID       int    `json:"id"`
	Title        string `json:"original_title"`
	TagLine      string `json:"tagline"`
	Description  string `json:"overview"`
	BackdropPath string `json:"backdrop_path"`
	PosterPath   string `json:"poster_path"`
	Year         string `json:"release_date"`
	Time         int    `json:"runtime"`
	Language     string `json:"original_language"`
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
	//encode search string
	endpoint, err := url.Parse(t.baseURL + "search/movie")
	if err != nil {
		return 0, fmt.Errorf("failed to fetch baseurl: %w", err)
	}

	params := url.Values{}
	params.Add("query", search)
	endpoint.RawQuery = params.Encode()

	// request
	req, err := http.NewRequest("GET", endpoint.String(), nil)
	fmt.Print(req)
	if err != nil {
		return 0, fmt.Errorf("Failed request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+t.apiKey)

	resp, err := t.client.Do(req)
	fmt.Println(resp.StatusCode)
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

func (t *TMDB) FetchFilmDetails(id int) (Films, error) {
	// url with query parameter safer approach me thinks
	endpoint, err := url.Parse(fmt.Sprintf("%s/movie/%d", t.baseURL, id))
	if err != nil {
		return Films{}, fmt.Errorf("invalid base url %w", err)
	}

	params := url.Values{}
	params.Add("append_to_response", "credits")
	endpoint.RawQuery = params.Encode()

	// start with the request
	req, err := http.NewRequest("GET", endpoint.String(), nil)
	if err != nil {
		return Films{}, errors.New("Request failed")
	}

	req.Header.Set("Authorization", "Bearer "+t.apiKey)
	resp, err := t.client.Do(req)
	if err != nil {
		return Films{}, fmt.Errorf("search failed: %w", err)
	}

	defer resp.Body.Close()
	var data TMDBFilmDetails

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return Films{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return Films{
		Title:          data.Title,
		TMDB_ID:        data.FilmID,
		Year:           data.Year,
		Poster:         data.PosterPath,
		BackdropPoster: data.BackdropPath,
		Description:    data.Description,
		TagLine:        data.TagLine,
		Time:           data.Time,
	}, nil

}
