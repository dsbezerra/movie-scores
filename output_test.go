package main

import (
	"io/ioutil"
	"os"
	"testing"
)

type TestOutputFormat struct {
	FieldOne   int      `json:"field_one"`
	FieldTwo   string   `json:"field_two"`
	FieldThree []string `json:"field_three"`
}

func TestOutputFile(t *testing.T) {
	data := TestOutputFormat{
		FieldOne:   1,
		FieldTwo:   "two",
		FieldThree: []string{"one", "two", "three"},
	}
	result := OutputFile("test_file", data)
	if result == nil {
		t.Errorf("Result was invalid, got: nil, expected a valid pointer\n")
	}

	contents, err := ioutil.ReadFile(result.Filename)
	if err != nil {
		t.Errorf("Failed to read file")
	}

	os.Remove(result.Filename)

	str := string(contents)
	if str != result.Data {
		t.Errorf("Contents was invalid, got: %s, expected: %s\n", str, result.Data)
	}
}
