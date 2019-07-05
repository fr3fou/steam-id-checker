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
// TODO think of what should this method return so that it can be used anywhere
func CheckIds(words io.Reader, key string, workerAmount int) {
	wordsScanner := bufio.NewScanner(words)
	var wg sync.WaitGroup

	for wordsScanner.Scan() {
		id := wordsScanner.Text()
		wg.Add(1)
		go checkID(id, key, &wg)
	}

	wg.Wait()
	fmt.Println("done")
}

func checkID(id, key string, wg *sync.WaitGroup) {
	// TODO: error handling
	resp, _ := steamapi.ResolveVanityURL(id, key)

	if !(resp.Response.Success == 1) {
		// Probs change to something else other than fmt.Println?
		fmt.Printf("%s is not taken on Steam!\n", id)
	} else {
		// Probs change to something else other than fmt.Println?
		fmt.Printf("%s is taken on Steam!\n", id)
	}

	wg.Done()
}
