package main

import (
	"bufio"
	"fmt"
	"log"
	"strconv"

	"flag"
	"io/ioutil"
	"os"
	"strings"

	"github.com/fr3fou/steam-id-checker/checker"
)

func main() {
	isInteractive := false
	flag.BoolVar(&isInteractive, "interactive", false, "display an interactive prompt to check IDs")
	flag.BoolVar(&isInteractive, "i", false, "display an interactive prompt to check IDs")

	filePath := ""
	flag.StringVar(&filePath, "file", "example", "path to the file which contains the IDs")
	flag.StringVar(&filePath, "f", "example", "path to the file which contains the IDs")

	workerAmount := 10
	flag.IntVar(&workerAmount, "workers", 10, "path to the file which contains the IDs")
	flag.IntVar(&workerAmount, "w", 10, "path to the file which contains the IDs")

	flag.Parse()

	if isInteractive {
		interactiveCli()
		return
	}

	finished := make(chan checker.SteamID)

	fi, err := os.Stdin.Stat()

	if err != nil {
		log.Fatal(err)
	}

	if fi.Mode()&os.ModeNamedPipe == 0 {
		if filePath == "" {
			log.Fatal("you need to pass in a file using the --file or -f flag")
		}

		file, err := os.Open(filePath)

		if err != nil {
			log.Fatal("you need to pass in a valid file using the --file or -f flag", err)
		}

		go checker.CheckIDs(file, workerAmount, finished)

		for val := range finished {
			if !val.IsTaken {
				fmt.Println(val.ID)
			}
		}

		return
	}

	go checker.CheckIDs(os.Stdin, workerAmount, finished)

	for val := range finished {
		if !val.IsTaken {
			fmt.Println(val.ID)
		}
	}

	// we have things coming through a pipe
}

func interactiveCli() {
	scanner := bufio.NewScanner(os.Stdin)
	finished := make(chan checker.SteamID)

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

	// we don't need the API key if the user is going to be scraping
	if checkForScrapingMethod(method) {
		go checker.CheckIDs(file, workerAmount, finished)
	} else {
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

	}

	f, err := os.OpenFile("results", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)

	if err != nil {
		log.Fatal(err)
	}

	for val := range finished {
		fmt.Println(val.Msg)

		if !val.IsTaken {
			_, err = f.WriteString(val.ID + "\n")
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	f.Close()
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
