package main

import "testing"

func TestImdbSearch(t *testing.T) {
	imdb := NewIMDb()

	query := "iron man 2008"

	result, err := imdb.Search(query)
	if err != nil {
		t.Error(err)
	}

	expectedQ := "iron_man_2008"
	expectedFirstItem := &imdbSearchItem{
		ID:      "tt0371746",
		Label:   "Iron Man",
		Subline: "Robert Downey Jr., Gwyneth Paltrow",
		Year:    2008,
	}

	if len(result.Data) == 0 {
		t.Errorf("Data was incorrect, got: 0, expected > 0")
	}

	if expectedQ != result.Query {
		t.Errorf("Query was incorrect, got: %s, expected: %s.", result.Query, query)
	}

	if !isSearchItemEqual(*expectedFirstItem, result.Data[0]) {
		t.Errorf("First item was incorrect, got: %s, expected: %s.", result.Data[0].ID, expectedFirstItem.ID)
	}
}

func isSearchItemEqual(a imdbSearchItem, b imdbSearchItem) bool {
	return (a.ID == b.ID &&
		a.Label == b.Label &&
		a.Subline == b.Subline &&
		a.Year == b.Year)
}
