package checker

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/fr3fou/go-steamapi"
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
// the reader to check whether the given ID exists
func CheckIDsWithAPI(words io.Reader, key string) ([]SteamID, error) {
	wordsScanner := bufio.NewScanner(words)

	ids := make([]SteamID, 0)
	for wordsScanner.Scan() {
		id := wordsScanner.Text()

		res, err := checkIDWithAPI(id, key)

		if err != nil {
			return nil, err
		}

		ids = append(ids, res)
	}

	err := wordsScanner.Err()

	if err != nil {
		return nil, err
	}

	return ids, nil
}

// CheckIDs takes in an io.Reader and scrapes the webpage against each word in
// the reader with workerAmount of workers to check whether the given ID exists
func CheckIDs(words io.Reader) ([]SteamID, error) {
	wordsScanner := bufio.NewScanner(words)

	ids := make([]SteamID, 0)
	for wordsScanner.Scan() {
		id := wordsScanner.Text()

		res, err := CheckID(id)

		if err != nil {
			return nil, err
		}

		ids = append(ids, res)
	}

	err := wordsScanner.Err()

	if err != nil {
		return nil, err
	}

	return ids, nil
}

func CheckID(id string) (SteamID, error) {
	// Sometimes the Steam servers give a false message, signifying that an ID is not taken
	// This occurs at random / if the severs are overloaded
	// We can work around this by making a few more requests and checking for its true value
	var (
		res SteamID
		err error
	)

	for i := 0; i < 3; i++ {
		res, err = checkID(id)

		if err != nil {
			return SteamID{}, err
		}

		// We don't need to make any more requests if the ID is taken (Steam servers sent the correct msg)
		if res.IsTaken {
			return res, nil
		}
	}

	return res, err
}

func checkID(id string) (SteamID, error) {
	url := fmt.Sprintf("http://steamcommunity.com/id/%s", id)
	resp, err := http.Get(url)

	if err != nil {
		return SteamID{}, err
	}

	defer resp.Body.Close()

	html, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return SteamID{}, err
	}

	if strings.Contains(string(html), "<h3>The specified profile could not be found.</h3>") {
		return SteamID{
			ID:      id,
			IsTaken: false,
			Msg:     fmt.Sprintf("%s is not taken on Steam!", id),
		}, nil
	}

	return SteamID{
		ID:      id,
		IsTaken: true,
		Msg:     fmt.Sprintf("%s is taken on Steam!", id),
	}, nil
}

func CheckIDWithAPI(id, key string) (SteamID, error) {
	// TODO: error handling
	resp, err := steamapi.ResolveVanityURL(id, key)

	if err != nil {
		return SteamID{}, err
	}

	if !(resp.Response.Success == 1) {
		return SteamID{
			ID:      id,
			IsTaken: false,
			Msg:     fmt.Sprintf("%s is not taken on Steam!", id),
		}, nil
	}

	return SteamID{
		ID:      id,
		IsTaken: true,
		Msg:     fmt.Sprintf("%s is taken on Steam!", id),
	}, nil
}
