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
	fmt.Println("HeMePi - Heavy Metal PI")

	flag.Var(&metal, "metal", "Metal to get price for")
	flag.Var(&currency, "currency", "Currency to get metal price in")

	flag.Parse()

	if *apiKey == "" {
		log.Fatalf("apikey must set, obtain one from goldapi.io")
	}

	client := NewGoldapiClient(baseURL, *apiKey, metal, currency)
	data, err := client.get()
	if err != nil {
		log.Fatal(err)
	}

	img, err := generateImage(data)
	if err != nil {
		log.Fatal(err)
	}

	err = render(img)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Thankyou for using HeMePi.")
}
