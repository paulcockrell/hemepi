package main

import (
	"flag"
	"fmt"
	"image/png"
	"log"
	"os"

	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/conn/spi/spireg"
	"periph.io/x/periph/experimental/devices/inky"
	"periph.io/x/periph/host"
)

var (
	spiPort  = flag.String("spi", "spidev0.0", "Name or number of SPI port to open")
	path     = flag.String("image", "", "Path to a png file to display on the inky")
	dcPin    = flag.String("dc", "22", "Inky DC pin")
	resetPin = flag.String("reset", "27", "Inky reset pin")
	busyPin  = flag.String("busy", "17", "Inky busy pin")

	model       = inky.PHAT
	modelColor  = inky.Red
	borderColor = inky.Black
)

func draw() error {
	flag.Var(&model, "model", "Inky model (PHAT or WHAT)")
	flag.Var(&modelColor, "model-color", "Inky model color (black, red or yellow)")
	flag.Var(&borderColor, "border-color", "Border color( black, white, red or yellow")
	flag.Parse()

	// Open and decode the image
	f, err := os.Open(*path)
	if err != nil {
		return err
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		return err
	}

	if _, err := host.Init(); err != nil {
		return err
	}

	log.Printf("Opening %s...", *spiPort)
	b, err := spireg.Open(*spiPort)
	if err != nil {
		return err
	}

	log.Printf("Opening pins...")
	dc := gpioreg.ByName(*dcPin)
	if dc == nil {
		return fmt.Errorf("invalid DC pin name: %s", *dcPin)
	}

	reset := gpioreg.ByName(*resetPin)
	if reset == nil {
		return fmt.Errorf("invalid Reset pin name: %s", *resetPin)
	}

	busy := gpioreg.ByName(*busyPin)
	if busy == nil {
		return fmt.Errorf("invalid Busy pin name: %s", *busyPin)
	}

	log.Printf("Creating inky...")
	dev, err := inky.New(b, dc, reset, busy, &inky.Opts{
		Model:       model,
		ModelColor:  modelColor,
		BorderColor: borderColor,
	})
	if err != nil {
		return err
	}

	log.Printf("Drawing image...")
	return dev.DrawAll(img)
}
