package checker

import (
	"bufio"
	"fmt"
	"io"

	"github.com/fr3fou/go-steamapi"
)

// CheckIds takes in an io.Reader and calls the Steam API against each word in
// the reader with workerAmont of workers to check whether the given ID exists
func CheckIds(words io.Reader, key string, workerAmount int) {
	fileScanner := bufio.NewScanner(words)

	for fileScanner.Scan() {
		id := fileScanner.Text()
		resp, err := steamapi.ResolveVanityURL(id, key)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(resp.Response.Success)
	}
}
