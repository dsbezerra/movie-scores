package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var opScore = "score"
var opSearch = "search"

var supportedOperations = []string{opScore, opSearch}
var supportedProviders = []string{IMDB, RottenT}

type (
	// Context represents the main application context
	Context struct {
		Provider  string
		Operation string
		Filename  string
		Query     string
		ID        string
	}

	// TODO: Make one result struct for both operations?

	// ScoreResult represents the result for a score operation
	ScoreResult struct {
		Provider   string  `json:"provider"`
		ID         string  `json:"id"`
		Score      float32 `json:"score"`
		ScoreClass string  `json:"score_class,omitempty"`
	}

	// SearchResult represents the result for a search operation
	SearchResult struct {
		Provider   string  `json:"provider"`
		ID         string  `json:"id"`
		Title      string  `json:"title"`
		Poster     string  `json:"poster"`
		Score      float32 `json:"score,omitempty"`
		ScoreClass string  `json:"score_class,omitempty"`
		Year       uint    `json:"year"`
	}

	// Provider is an interface used to reduce equal code
	Provider interface {
		Score(id string) (*ScoreResult, error)
		Search(query string) ([]SearchResult, error)
	}
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
		return nil, errors.New("response code was not successfull")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func isArgValid(arg string, collection []string) bool {
	for _, i := range collection {
		if i == arg {
			return true
		}
	}
	return false
}

func isOpeprationSupported(op string) bool {
	return isArgValid(op, supportedOperations)
}

func isProviderSupported(provider string) bool {
	return isArgValid(provider, supportedProviders)
}

func checkArgs() *Context {
	/**
	 * -p [Required]
	 * Provider used in operation.
	 *
	 * imdb   - IMDb: https://imdb.com.br/
	 * rotten - RottenTomatoes: https://www.rottentomatoes.com/
	 */
	provider := flag.String("p", "", "Provider to process (imdb/rotten)")

	/**
	* -op [Required]
	* Operation to run.
	*
	* search - Uses provider's default search API to search for movies. Returns a list as result.
	* score  - Uses given ID to retrieve movie score.
	 */
	operation := flag.String("op", "", "Operation to execute (search/score)")

	/**
	* -out [Required]
	* Filename of the outputted file with results.
	 */
	filename := flag.String("out", "", "Filename to output")

	/**
	* -q [Required if operation is search]
	* Query used in search operations.
	 */
	query := flag.String("q", "", "Query used in search operations")

	/**
	* -id [Required if operation is score]
	* Identifier used in score operations.
	 */
	id := flag.String("id", "", "Identifier used in score operations")

	flag.Parse()

	if *provider == "" || *operation == "" || *filename == "" {
		log.Fatalf("Error: all parameters must be defined")
	}

	if !isProviderSupported(*provider) {
		log.Fatalf("Error: provider '%s' is not supported", *provider)
	}

	if !isOpeprationSupported(*operation) {
		log.Fatalf("Error: operation '%s' is not supported", *operation)
	}

	if *operation == "search" && *query == "" {
		log.Fatalf("Error: query is required for search operation")
	}

	if *operation == "score" && *id == "" {
		log.Fatalf("Error: id is required for score operation")
	}

	return &Context{
		Provider:  *provider,
		Filename:  *filename,
		Operation: *operation,
		Query:     *query,
		ID:        *id,
	}
}

func (ctx *Context) run() {
	var p Provider

	if ctx.Provider == IMDB {
		p = NewIMDb()
	} else if ctx.Provider == RottenT {
		p = NewRottenTomatoes()
	}

	if p != nil {
		var result interface{}
		var err error

		switch ctx.Operation {
		case opSearch:
			result, err = p.Search(ctx.Query)
			if err != nil {
				log.Fatal(err)
			}
		case opScore:
			result, err = p.Score(ctx.ID)
			if err != nil {
				log.Fatal(err)
			}
		default:
		}

		if result != nil {
			r := OutputFile(ctx.Filename, result)
			fmt.Printf("Outputted to: %s\n", r.Filename)
		}
	}
}

func main() {
	ctx := checkArgs()
	ctx.run()
}
