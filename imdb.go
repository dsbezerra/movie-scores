package main

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
)

const IMDB = "imdb"

const imdbApiBaseURL = "https://v2.sg.media-imdb.com/suggests/"

type (
	IMDb struct {
	}

	imdbSearchResult struct {
		V     int              `json:"v"`
		Query string           `json:"q"`
		Data  []imdbSearchItem `json:"d"`
	}

	imdbSearchItem struct {
		ID      string                `json:"id"`
		Image   []interface{}         `json:"i"`
		Label   string                `json:"l"`
		Q       string                `json:"q"`
		Subline string                `json:"s"`
		VT      int                   `json:"vt"`
		Videos  []imdbSearchItemVideo `json:"v"`
		Year    uint                  `json:"y"`
	}

	imdbSearchItemVideo struct {
		Images  []interface{} `json:"i"`
		ID      string        `json:"id"`
		Label   string        `json:"l"`
		Subline string        `json:"s"`
	}

	imdbSearchOutputFormat struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Poster string `json:"poster"`
	}
)

func NewIMDb() *IMDb {
	return &IMDb{}
}

// Search returns movies for a given query from IMDB suggests API
func (imdb *IMDb) Search(query string) ([]SearchResult, error) {
	if query == "" {
		return nil, nil
	}

	fullURL := imdbApiBaseURL + string(query[0]) + "/"

	query = strings.Replace(query, " ", "_", -1)
	fullURL += query + ".json"

	body, err := Get(fullURL)
	if err != nil {
		return nil, err
	}

	str := string(body)
	start := strings.Index(str, "(") + 1
	end := strings.LastIndex(str, ")")
	if start > 0 && end > start {
		str = str[start:end]
	} else {
		return nil, errors.New("couldn't find string between `(` `)` tokens")
	}

	var result imdbSearchResult
	err = json.Unmarshal([]byte(str), &result)
	if err != nil {
		return nil, err
	}

	r := make([]SearchResult, 0)
	for _, item := range result.Data {
		sr := SearchResult{
			Provider: IMDB,
			ID:       item.ID,
			Title:    item.Label,
			Year:     item.Year,
		}

		if len(item.Image) > 0 {
			sr.Poster = getString(item.Image[0])
		}

		r = append(r, sr)
	}

	return r, nil
}

func (imdb *IMDb) Score(id string) (*ScoreResult, error) {
	return nil, nil
}

func (searchResult *imdbSearchResult) toOutputFormat() []imdbSearchOutputFormat {
	result := make([]imdbSearchOutputFormat, 0)

	for _, item := range searchResult.Data {
		r := imdbSearchOutputFormat{
			ID:   item.ID,
			Name: item.Label,
		}

		if len(item.Image) > 0 {
			r.Poster = getString(item.Image[0])
		}

		result = append(result, r)
	}
	return result
}

func getString(val interface{}) string {
	typ := reflect.TypeOf(val)
	if typ != nil && typ.Kind() == reflect.String {
		return reflect.ValueOf(val).String()
	}
	return ""
}
