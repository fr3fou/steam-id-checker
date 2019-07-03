package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/fr3fou/go-steamapi"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter the path to txt file: ")

	scanner.Scan()
	path := scanner.Text()

	fmt.Print("Enter enter your Steam API key (you can get yours at https://steamcommunity.com/dev/apikey): ")

	scanner.Scan()
	key := scanner.Text()

	file, err := os.Open(path)

	if err != nil {
		log.Fatal("The path you entered isn't a valid file!", err)
	}

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
