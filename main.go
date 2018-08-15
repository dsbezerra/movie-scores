package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// Get performs a GET request to the given URL
func Get(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.106 Safari/537.36")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	if response.StatusCode != 200 {
		// TODO: better handle this
		return nil, errors.New("response code was not sucessfull")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func main() {

	rotten := NewRottenTomatoes()

	query := "iron man"

	result, err := rotten.Search(query)
	if err != nil {
		log.Fatal(err)
	}

	data := result.toOutputFormat()
	outRes := OutputToFile("result_rotten_", data)

	fmt.Printf("File outputted to: %s\nContents: %s\n", outRes.Filename, outRes.Data)
	defer os.Remove(outRes.Filename)
}
