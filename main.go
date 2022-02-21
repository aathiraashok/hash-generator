package main

import (
	"crypto/md5"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const DefaultparallelRequests = 10

var parallelRequests = flag.Int("parallel", DefaultparallelRequests, "Parallel requests limit")

var usage = `Usage: ./myhttp [options...] [urls...]

Options:
	-parallel Number of workers to run concurrently. Default is 10.
`

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, fmt.Sprintf(usage))
	}
	flag.Parse()
	urls := flag.Args()
	workersCount := *parallelRequests
	if workersCount > len(urls) {
		workersCount = len(urls)
	}
	tool := Tool{
		Parallel: workersCount,
	}
	x := tool.Run(urls...)
	fmt.Println(x)
}

type Tool struct {
	Parallel int
}

func (t *Tool) Run(urls ...string) error {
	var validURLs []string

	for _, u := range urls {

		if x, err := url.Parse(u); err != nil {
			return fmt.Errorf("could not parse %s, error: %v", u, err)
		} else {
			validURLs = append(validURLs, x.String())
		}
	}
	urlChannel := make(chan string, len(urls))

	// Finished channel will be used to be sure that the worker is finihed its job
	finished := make(chan bool, len(urls))
	// Create workers
	for i := 0; i < t.Parallel; i++ {
		go t.worker(urlChannel, finished)
	}

	// Send URLs
	for _, url := range urls {
		urlChannel <- url
	}

	close(urlChannel)
	for i := 0; i < len(urls); i++ {
		<-finished
	}
	close(finished)
	return nil
}

func (t Tool) worker(urlChannel <-chan string, finished chan<- bool) {
	// Whenever channel sends a URL, one of worker get it
	for url := range urlChannel {
		// Send HTTP request, read the body and create hash
		hash, err := createHash(url)

		// // If there is an error in the request, print it
		if err != nil {
			fmt.Println(url, err.Error())
			finished <- true
			continue
		}
		fmt.Println(url, hash)
		finished <- true
	}
}
func createHash(url string) (string, error) {
	if !strings.HasPrefix(url, "http") {
		url = "http://" + url
	}
	client := &http.Client{}
	resp, e := client.Get(url)
	if e != nil {
		return "", e
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}

	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return "", e
	}
	val := (md5.Sum(body))
	return fmt.Sprintf("%x", val), nil
}
