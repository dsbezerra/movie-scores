package main

import "testing"

func TestImdbSearch(t *testing.T) {
	imdb := NewIMDb()

	query := "iron man 2008"

	result, err := imdb.Search(query)
	if err != nil {
		t.Error(err)
	}

	expectedFirstItem := &imdbSearchItem{
		ID:      "tt0371746",
		Label:   "Iron Man",
		Subline: "Robert Downey Jr., Gwyneth Paltrow",
		Year:    2008,
	}

	if len(result) == 0 {
		t.Errorf("Size was incorrect, got: 0, expected > 0")
	}

	if !isSearchItemEqual(result[0], *expectedFirstItem) {
		t.Errorf("First item was incorrect, got: %s, expected: %s.", result[0].ID, expectedFirstItem.ID)
	}
}

func isSearchItemEqual(a SearchResult, b imdbSearchItem) bool {
	return a.ID == b.ID
}
