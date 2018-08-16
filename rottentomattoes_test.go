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
		t.Errorf("Movie was not in list.\n")
	}
}

func isMovieEqual(a SearchResult, b rtMovie) bool {
	return (a.Title == b.Name &&
		a.ID == b.URL &&
		a.Year == b.Year)
}
