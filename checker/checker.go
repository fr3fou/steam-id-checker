package checker

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/fr3fou/go-steamapi"
	"github.com/ti/nasync"
)

// SteamID is the underlying struct
// which is used to communicate using the channels between the goroutines
// or the caller of CheckIDs / CheckIDsWithAPI
type SteamID struct {
	ID      string
	IsTaken bool
	Msg     string
}

// CheckIDsWithAPI takes in an io.Reader and calls the Steam API against each word in
// the reader with workerAmont of workers to check whether the given ID exists
func CheckIDsWithAPI(words io.Reader, key string, workerAmount int, finished chan SteamID) {
	async := nasync.New(workerAmount, workerAmount)
	defer async.Close()

	wordsScanner := bufio.NewScanner(words)
	var wg sync.WaitGroup

	for wordsScanner.Scan() {
		id := wordsScanner.Text()
		wg.Add(1)

		async.Do(func() {
			checkIDWithAPI(id, key, &wg, finished)
		})
	}

	wg.Wait()
	close(finished)
}

// CheckIDs takes in an io.Reader and scrapes the webpage against each word in
// the reader with workerAmount of workers to check whther the given ID exists
func CheckIDs(words io.Reader, workerAmount int, finished chan SteamID) {
	async := nasync.New(workerAmount, workerAmount)
	defer async.Close()

	wordsScanner := bufio.NewScanner(words)
	var wg sync.WaitGroup

	for wordsScanner.Scan() {
		id := wordsScanner.Text()
		wg.Add(1)

		async.Do(func() {
			checkID(id, &wg, finished)
		})
	}

	wg.Wait()
	close(finished)
}

func checkID(id string, wg *sync.WaitGroup, finished chan SteamID) {
	url := fmt.Sprintf("http://steamcommunity.com/id/%s", id)
	resp, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	defer wg.Done()

	html, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	if strings.Contains(string(html), "<h3>The specified profile could not be found.</h3>") {
		finished <- SteamID{
			ID:      id,
			IsTaken: false,
			Msg:     fmt.Sprintf("%s is not taken on Steam!", id),
		}
	} else {
		finished <- SteamID{
			ID:      id,
			IsTaken: true,
			Msg:     fmt.Sprintf("%s is taken on Steam!", id),
		}
	}
}

func checkIDWithAPI(id, key string, wg *sync.WaitGroup, finished chan SteamID) {
	// TODO: error handling
	resp, _ := steamapi.ResolveVanityURL(id, key)
	defer wg.Done()

	if !(resp.Response.Success == 1) {
		finished <- SteamID{
			ID:      id,
			IsTaken: false,
			Msg:     fmt.Sprintf("%s is not taken on Steam!", id),
		}
	} else {
		finished <- SteamID{
			ID:      id,
			IsTaken: true,
			Msg:     fmt.Sprintf("%s is taken on Steam!", id),
		}
	}
}
