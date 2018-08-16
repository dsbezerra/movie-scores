package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"unicode"

	"github.com/PuerkitoBio/goquery"
)

const ROTTEN_TOMATOES = "rotten"

const rottenBaseURL = "https://www.rottentomatoes.com"
const rottenApiBaseURL = rottenBaseURL + "/napi/"

const CLASS_ROTTEN = "rotten"
const CLASS_FRESH = "fresh"
const CLASS_CERTIFIED_FRESH = "certified_fresh"

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

	rtScoreResult struct {
		Name       string
		MeterClass string
		MeterScore uint
		Path       string
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
func (rt *RottenTomatoes) Search(query string) ([]SearchResult, error) {
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

	r := make([]SearchResult, 0)
	for _, movie := range result.Movies {
		r = append(r, SearchResult{
			ID:         movie.URL,
			Title:      movie.Name,
			Poster:     movie.Image,
			Provider:   ROTTEN_TOMATOES,
			Score:      float32(movie.MeterScore),
			ScoreClass: movie.MeterClass,
			Year:       movie.Year,
		})
	}

	return r, nil
}

func (rt *RottenTomatoes) Score(id string) (*ScoreResult, error) {
	if id == "" {
		return nil, nil
	}

	path := id
	finalPath := ensurePathHasM(path)
	fullURL := rottenBaseURL + finalPath
	body, err := Get(fullURL)
	if err != nil {
		return nil, err
	}

	result := &rtScoreResult{}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	container := doc.Find("#all-critics-numbers > div > div:nth-child(1) > div > div.critic-score.meter")

	result.Name = strings.TrimSpace(doc.Find("#heroImageContainer > a > h1").Text())
	result.Path = finalPath

	meterScoreText := container.Find("span.meter-value.superPageFontColor > span").Text()
	result.MeterScore = scoreAsInt(meterScoreText)

	icon := container.Find("span.meter-tomato.icon")
	val, exists := icon.Attr("class")
	if exists {
		if strings.Contains(val, CLASS_ROTTEN) {
			result.MeterClass = CLASS_ROTTEN
		} else if strings.Contains(val, CLASS_FRESH) {
			result.MeterClass = CLASS_FRESH
		} else if strings.Contains(val, CLASS_CERTIFIED_FRESH) {
			result.MeterClass = CLASS_CERTIFIED_FRESH
		}
	}

	return &ScoreResult{
		Provider:   ROTTEN_TOMATOES,
		ID:         result.Path,
		Score:      float32(result.MeterScore),
		ScoreClass: result.MeterClass,
	}, nil
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

func ensurePathHasM(path string) string {
	if strings.HasPrefix(path, "/m/") {
		return path
	} else if !strings.HasPrefix(path, "/m/") && strings.HasPrefix(path, "m/") {
		path = "/" + path
	} else if !strings.HasPrefix(path, "m/") {
		if path[0] == '/' {
			path = "/m" + path
		} else {
			path = "/m/" + path
		}
	}
	return path
}

func scoreAsInt(scoreText string) uint {
	var result uint
	trimmed := strings.TrimFunc(scoreText, func(r rune) bool {
		return !unicode.IsNumber(r)
	})

	if trimmed != "" {
		number, err := strconv.Atoi(trimmed)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
		}

		result = uint(number)
	}

	return result
}
