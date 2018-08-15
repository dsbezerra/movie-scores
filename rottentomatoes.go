package main

import (
	"encoding/json"
	"net/url"
)

const rottenApiBaseURL = "https://www.rottentomatoes.com/napi/"

type (
	RottenTomatoes struct {
	}

	/* Response struct for url:
	   https://www.rottentomatoes.com/napi/search?query="something"
	*/
	rtSearchResult struct {
		ActorCount     uint          `json:"actorCount"`
		Actors         []rtActor     `json:"actors"`
		CriticCount    uint          `json:"criticCount"`
		Critics        []rtCritic    `json:"critics"`
		FranchiseCount uint          `json:"franchiseCount"`
		Franchises     []rtFranchise `json:"franchises"`
		MovieCount     uint          `json:"movieCount"`
		Movies         []rtMovie     `json:"movies"`
		TvCount        uint          `json:"tvCount"`
		TvSeries       []rtTvShow    `json:"tvSeries"`
	}

	rtActor struct {
		Image string `json:"image"`
		Name  string `json:"name"`
		URL   string `json:"url"`
	}

	rtCritic struct {
		Image        string   `json:"image"`
		Name         string   `json:"name"`
		Publications []string `json:"publications"`
		URL          string   `json:"url"`
	}

	rtFranchise struct {
		Image string `json:"image"`
		Title string `json:"title"`
		URL   string `json:"url"`
	}

	rtMovie struct {
		CastItems  []rtCastItem `json:"castItems"`
		Image      string       `json:"image"`
		MeterClass string       `json:"meterClass"`
		MeterScore uint         `json:"meterScore"`
		Name       string       `json:"name"`
		Subline    string       `json:"subline"`
		URL        string       `json:"url"`
		Year       uint         `json:"year"`
	}

	rtCastItem struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	rtTvShow struct {
		Image      string `json:"image"`
		MeterClass string `json:"meterClass"`
		MeterScore uint   `json:"meterScore"`
		StartYear  uint   `json:"startYear"`
		EndYear    uint   `json:"endYear"`
		Title      string `json:"title"`
		URL        string `json:"url"`
	}

	rtSearchOutputFormat struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		Poster     string `json:"poster"`
		MeterClass string `json:"meterClass"`
		MeterScore uint   `json:"meterScore"`
	}
)

// NewRottenTomatoes Gives an instance of rotten struct
func NewRottenTomatoes() *RottenTomatoes {
	return &RottenTomatoes{}
}

// Search for movie, actors, shows, franchises, etc, using rotten public api.
func (rt *RottenTomatoes) Search(query string) (*rtSearchResult, error) {
	if query == "" {
		return nil, nil
	}

	query = url.QueryEscape(query)
	url := rottenApiBaseURL + "search/?limit=5&query=" + query

	body, err := Get(url)
	if err != nil {
		return nil, err
	}

	var result rtSearchResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (rt *rtSearchResult) toOutputFormat() []rtSearchOutputFormat {
	result := make([]rtSearchOutputFormat, 0)

	// NOTE: only support movies
	if rt.MovieCount == 0 {
		return result
	}

	for _, movie := range rt.Movies {
		m := rtSearchOutputFormat{
			ID:         movie.URL,
			Name:       movie.Name,
			Poster:     movie.Image,
			MeterClass: movie.MeterClass,
			MeterScore: movie.MeterScore,
		}

		result = append(result, m)
	}

	return result
}
