package checker

import (
	"bufio"
	"io"

	"github.com/fr3fou/go-steamapi"
)

type idChecker struct {
	id  string
	key string
}

// CheckIds takes in an io.Reader and calls the Steam API against each word in
// the reader with workerAmont of workers to check whether the given ID exists
func CheckIds(words io.Reader, key string, workerAmount int) {
	fileScanner := bufio.NewScanner(words)

	jobs := make(chan idChecker)
	results := make(chan *idChecker)

	// Start up workerAmount of workers
	for w := 1; w < workerAmount; w++ {
		go worker(jobs, results)
	}

	// Fill in the jobs queue (channel)
	for fileScanner.Scan() {
		id := fileScanner.Text()
		jobs <- idChecker{
			id,
			key,
		}
	}

	close(jobs)
}

func worker(jobs <-chan idChecker, results chan<- *idChecker) {
	// for every job in the queue (channel), call checkID
	for j := range jobs {
		results <- checkID(j)
	}
}

func checkID(ic idChecker) *idChecker {
	resp, err := steamapi.ResolveVanityURL(ic.id, ic.key)

	if err != nil {
		return nil
	}

	// According to the documentation, 42 is the code returned when an ID doesn't exist
	if resp.Response.Success == 42 {
		return &ic
	}

	return nil
}
