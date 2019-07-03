package main

import (
	"bufio"
	"fmt"
	"log"
	"strconv"

	"io/ioutil"
	"os"
	"strings"

	"github.com/fr3fou/go-steamapi"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter the path to txt file: ")

	scanner.Scan()
	path := scanner.Text()

	file, err := os.Open(path)

	if err != nil {
		log.Fatal("The path you entered isn't a valid file!", err)
	}

	key := ""

	contents, err := ioutil.ReadFile(".key")

	// If we have read the file successfully, this means that there must be a key inside
	if err == nil {
		fmt.Printf("An existing Steam API key has been found (%s...), would you like to use it? (Y/n): ", string(contents)[:6])
		scanner.Scan()
		useExisting := scanner.Text()

		if checkForAgreement(useExisting) {
			key = string(contents)
		}
	}

	// If the key still empty, ask for it
	if key == "" {
		fmt.Print("Enter enter your Steam API key (you can get yours at https://steamcommunity.com/dev/apikey): ")

		scanner.Scan()
		key = scanner.Text()

		fmt.Print("Would you like to remember the key for future use? (Y/n): ")

		scanner.Scan()
		remember := strings.Trim(scanner.Text(), " ")

		if checkForAgreement(remember) {
			ioutil.WriteFile(".key", []byte(key), 0644)
		}
	}

	fmt.Print("How many workers would you like to use? - how many IDs can be processed at a time (default is 50): ")

	scanner.Scan()
	workerAmount := 50
	workerInput := scanner.Text()

	if workerInput != "" {
		workerAmount, err = strconv.Atoi(workerInput)
		if err != nil {
			log.Fatal(err)
		}
	}

	checkIds(file, key, workerAmount)
}

func checkIds(file *os.File, key string, workerAmount int) {
	fileScanner := bufio.NewScanner(file)

	for fileScanner.Scan() {
		id := fileScanner.Text()
		resp, err := steamapi.ResolveVanityURL(id, key)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(resp.Response.Success)
	}
}

func checkForAgreement(s string) bool {
	return s == "y" || s == "Y" || s == "Yes" || s == "yes" || s == ""
}
