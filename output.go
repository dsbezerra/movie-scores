package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type (
	OutputResult struct {
		Filename string
		Data     string
	}
)

// OutputToFile outputs struct data to a JSON file.
func OutputToFile(filename string, data interface{}) *OutputResult {
	if filename == "" || data == nil {
		return nil
	}

	contents, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return nil
	}

	tmpfile, err := ioutil.TempFile("./data", filename)
	if _, err := tmpfile.Write(contents); err != nil {
		return nil
	}
	if err := tmpfile.Close(); err != nil {
		return nil
	}
	// NOTE: Don't forget to remove file from caller after process.
	// defer os.Remove(tmpfile.Name())
	return &OutputResult{
		Filename: tmpfile.Name(),
		Data:     string(contents),
	}
}
