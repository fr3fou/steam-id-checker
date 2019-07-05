package main

import (
	"bufio"
	"fmt"
	"log"
	"strconv"

	"io/ioutil"
	"os"
	"strings"

	"github.com/fr3fou/steam-id-checker/checker"
)

func main() {
	// TODO: check for interactive flag (-i)
	interactiveCli()
}

func interactiveCli() {
	scanner := bufio.NewScanner(os.Stdin)
	finished := make(chan string)

	fmt.Print("Enter the path to txt file: ")
	scanner.Scan()
	path := scanner.Text()

	file, err := os.Open(path)

	if err != nil {
		log.Fatal("The path you entered isn't a valid file!", err)
	}

	fmt.Print("Would you like to scrape or use the Steam API? (default is scrape): ")
	scanner.Scan()
	method := scanner.Text()

	fmt.Print("How many workers would you like to use? - how many IDs can be processed at a time (default is 10): ")

	scanner.Scan()
	workerAmount := 10
	workerInput := scanner.Text()

	if workerInput != "" {
		workerAmount, err = strconv.Atoi(workerInput)
		if err != nil {
			log.Fatal(err)
		}
	}

	if checkForScrapingMethod(method) {
		go checker.CheckIDs(file, workerAmount, finished)

		for val := range finished {
			fmt.Println(val)
		}

		return
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

	go checker.CheckIDsWithAPI(file, key, workerAmount, finished)

	for val := range finished {
		fmt.Println(val)
	}

	close(finished)
}

func checkForAgreement(s string) bool {
	return checkForAnswer(s, []string{"y", "Y", "Yes", "yes", ""})
}

func checkForScrapingMethod(s string) bool {
	return checkForAnswer(s, []string{"scrape", "web", ""})
}

func checkForAnswer(s string, answers []string) bool {
	for _, val := range answers {
		if val == s {
			return true
		}
	}

	return false
}
