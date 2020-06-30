package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
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

// UserAgents is a list of user agents that the scraper can use to trick the webmasters
// and hopefully don't get block
var UserAgents = [...]string{
	// Linus Firefox
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:43.0) Gecko/20100101 Firefox/43.0",
	// Mac Firefox
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.11; rv:43.0) Gecko/20100101 Firefox/43.0",
	// Mac Safari 4
	"Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_6_2; de-at) AppleWebKit/531.21.8 (KHTML, like Gecko) Version/4.0.4 Safari/531.21.10",
	// Mac Safari
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/601.3.9 (KHTML, like Gecko) Version/9.0.2 Safari/601.3.9",
	// Windows Chrome
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.125 Safari/537.36",
	// Windows IE 10
	"Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.2; WOW64; Trident/6.0)",
	// Windows IE 11
	"Mozilla/5.0 (Windows NT 6.3; WOW64; Trident/7.0; rv:11.0) like Gecko",
	// Windows Edge
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/46.0.2486.0 Safari/537.36 Edge/13.10586",
	// Windows Firefox
	"Mozilla/5.0 (Windows NT 6.3; WOW64; rv:43.0) Gecko/20100101 Firefox/43.0",
	// iPhone
	"Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B5110e Safari/601.1",
	// iPad
	"Mozilla/5.0 (iPad; CPU OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1",
	// Android
	"Mozilla/5.0 (Linux; Android 5.1.1; Nexus 7 Build/LMY47V) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.76 Safari/537.36",
}

// GetRandomUserAgent retrieves a random user agent
func GetRandomUserAgent() string {
	result := ""

	// Using current time nanosecond as seed
	seed := time.Now().Nanosecond()

	// Seed the random
	rand.Seed(int64(seed))

	// Get random user-agent
	size := len(UserAgents)
	result = UserAgents[rand.Int31n(int32(size))]

	return result
}

// Get performs a GET request to the given URL
func Get(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", GetRandomUserAgent())

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
	fmt.Println(response.Status)
	fmt.Println(string(body))

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

func isOperationSupported(op string) bool {
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

	if !isOperationSupported(*operation) {
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
