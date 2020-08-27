package main

import (
	"flag"
	"fmt"

	"log"

	"periph.io/x/periph/experimental/devices/inky"
)

var (
	baseURL = "https://www.goldapi.io/api"

	// Client flags
	apiKey = flag.String("apikey", "", "API key for goldapi.io")

	// InkyPHAT display flags
	spiPort  = flag.String("spi", "/dev/spidev0.0", "Name or number of SPI port to open")
	dcPin    = flag.String("dc", "22", "Inky DC pin")
	resetPin = flag.String("reset", "27", "Inky reset pin")
	busyPin  = flag.String("busy", "17", "Inky busy pin")

	metal    Metal
	currency Currency

	model       = inky.PHAT
	modelColor  = inky.Red
	borderColor = inky.Black
)

func main() {
	fmt.Print("HeMePI - [HE]vy [ME]tal Raspberry [PI] Gold and Silver price tracker\n\n")

	flag.Var(&metal, "metal", "Metal to get price for")
	flag.Var(&currency, "currency", "Currency to get metal price in")
	flag.Var(&model, "model", "Inky model (PHAT or WHAT)")
	flag.Var(&modelColor, "model-color", "Inky model color (black, red or yellow)")
	flag.Var(&borderColor, "border-color", "Border color( black, white, red or yellow")

	flag.Parse()

	if *apiKey == "" {
		log.Fatalf("apikey must set, obtain one from goldapi.io")
	}

	// Setup goldapi.io client and get currency pair
	client := NewGoldapiClient(baseURL, *apiKey, metal, currency)
	data, err := client.get()
	if err != nil {
		log.Fatal(err)
	}

	// Generate currency pair image
	img, err := generateImage(data)
	if err != nil {
		log.Fatal(err)
	}

	// Display image
	display, err := NewInky(
		spiPort,
		dcPin,
		resetPin,
		busyPin,
		&inky.Opts{
			Model:       model,
			ModelColor:  modelColor,
			BorderColor: borderColor,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	if err = display.DrawAll(*img); err != nil {
		log.Fatal(err)
	}

	log.Print("\n\nThankyou for using HeMePi.\n")
}
