package checker

import (
	"bufio"
	"fmt"
	"io"
	"sync"

	"github.com/fr3fou/go-steamapi"
)

// CheckIds takes in an io.Reader and calls the Steam API against each word in
// the reader with workerAmont of workers to check whether the given ID exists
func CheckIds(words io.Reader, key string, workerAmount int, finished chan string) {
	wordsScanner := bufio.NewScanner(words)
	var wg sync.WaitGroup

	for wordsScanner.Scan() {
		id := wordsScanner.Text()
		wg.Add(1)
		go checkID(id, key, &wg, finished)
	}

	wg.Wait()
	close(finished)
}

func checkID(id, key string, wg *sync.WaitGroup, finished chan string) {
	// TODO: error handling
	resp, _ := steamapi.ResolveVanityURL(id, key)
	defer wg.Done()

	if !(resp.Response.Success == 1) {
		finished <- fmt.Sprintf("%s is not taken on Steam!", id)
	} else {
		finished <- fmt.Sprintf("%s is taken on Steam!", id)
	}
}
