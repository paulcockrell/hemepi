package main

import (
	"flag"
	"fmt"

	"log"
)

var (
	baseURL = "https://www.goldapi.io/api"

	apiKey = flag.String("apikey", "", "API key for goldapi.io")

	metal    Metal
	currency Currency
)

func main() {
	flag.Var(&metal, "metal", "Metal to get price for")
	flag.Var(&currency, "currency", "Currency to get metal price in")

	flag.Parse()

	client := NewGoldapiClient(baseURL, *apiKey, metal, currency)
	data, err := client.get()
	if err != nil {
		log.Fatal(err)
	}

	err = generateImage()
	if err != nil {
		fmt.Println(err)
	}

	// For now just dump data to terminal
	fmt.Println(data)
}
