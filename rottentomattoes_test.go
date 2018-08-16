package main

import "testing"

func TestRottenSearch(t *testing.T) {
	rotten := NewRottenTomatoes()

	query := "iron man"
	result, err := rotten.Search(query)
	if err != nil {
		t.Error(err)
	}

	expectedMovieInList := rtMovie{
		Name: "Iron Man",
		Year: 2008,
		URL:  "/m/iron_man",
	}

	if len(result) == 0 {
		t.Errorf("Movie count was incorrect, got: 0, expected > 0")
	}

	expectedFound := false
	for _, movie := range result {
		expectedFound = isMovieEqual(movie, expectedMovieInList)
		if expectedFound {
			break
		}
	}

	if !expectedFound {
		t.Errorf("Movie was not in list.")
	}
}

func TestRottenScore(t *testing.T) {
	rotten := NewRottenTomatoes()

	// NOTE: this can break if movie score changes...
	path := "/m/sharknado_2013"
	result, err := rotten.Score(path)
	if err != nil {
		t.Error(err)
	}

	expectedScore := float32(82)
	expectedClass := "fresh"

	if expectedScore != result.Score {
		t.Errorf("Score was incorrect, got %f, expected: %f", result.Score, expectedScore)
	}

	if !isScoreClassOneOf(result.ScoreClass) {
		t.Errorf("Score class was incorrect, got %s, expected: %s", result.ScoreClass, expectedClass)
	}
}

func isMovieEqual(a SearchResult, b rtMovie) bool {
	return (a.Title == b.Name &&
		a.ID == b.URL &&
		a.Year == b.Year)
}

func isScoreClassOneOf(scoreClass string) bool {
	return scoreClass == "rotten" || scoreClass == "fresh" || scoreClass == "certified_fresh"
}
