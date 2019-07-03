package checker

import (
	"bufio"
	"fmt"
	"io"

	"github.com/fr3fou/go-steamapi"
)

type idChecker struct {
	id      string
	key     string
	isTaken bool
}

// CheckIds takes in an io.Reader and calls the Steam API against each word in
// the reader with workerAmont of workers to check whether the given ID exists
// TODO think of what should this method return so that it can be used anywhere
func CheckIds(words io.Reader, key string, workerAmount int) {
	wordsScanner := bufio.NewScanner(words)

	// TODO what should the buffer size be?
	// I think it has to be the total amount of words (but how do I get them from the reader?)
	jobs := make(chan idChecker)
	results := make(chan idChecker)

	// Start up workerAmount of workers
	for w := 1; w < workerAmount; w++ {
		go worker(jobs, results)
	}

	limit := 1

	// Fill in the jobs queue (channel)
	for ; wordsScanner.Scan(); limit++ {
		id := wordsScanner.Text()
		jobs <- idChecker{
			id:      id,
			key:     key,
			isTaken: false,
		}
	}

	close(jobs)

	for i := 0; i < limit; i++ {
		result := <-results
		if !result.isTaken {
			// Probs change to something else other than fmt.Println?
			fmt.Printf("(%d out of %d) - %s is not taken on Steam!", i, limit, result.id)
		} else {
			// Probs change to something else other than fmt.Println?
			fmt.Printf("(%d out of %d) - %s is taken on Steam.", i, limit, result.id)
		}
	}

	close(results)
}

func worker(jobs <-chan idChecker, results chan<- idChecker) {
	// for every job in the queue (channel), call checkID
	for j := range jobs {
		results <- checkID(j)
	}
}

func checkID(ic idChecker) idChecker {
	// TODO: error handling
	resp, _ := steamapi.ResolveVanityURL(ic.id, ic.key)
	ic.isTaken = resp.Response.Success == 1
	return ic
}
