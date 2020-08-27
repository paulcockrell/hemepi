package main

import (
	"fmt"
	"log"

	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/conn/spi/spireg"
	"periph.io/x/periph/experimental/devices/inky"
	"periph.io/x/periph/host"
)

// NewInky constructs new Inky display device after attempting to setup hardware
func NewInky(spiPort, dcPin, resetPin, busyPin *string, opts *inky.Opts) (*inky.Dev, error) {
	if _, err := host.Init(); err != nil {
		return nil, err
	}

	log.Printf("Opening %s...", *spiPort)
	b, err := spireg.Open(*spiPort)
	if err != nil {
		return nil, err
	}

	log.Printf("Opening pins...")
	dc := gpioreg.ByName(*dcPin)
	if dc == nil {
		return nil, fmt.Errorf("invalid DC pin name: %s", *dcPin)
	}

	reset := gpioreg.ByName(*resetPin)
	if reset == nil {
		return nil, fmt.Errorf("invalid Reset pin name: %s", *resetPin)
	}

	busy := gpioreg.ByName(*busyPin)
	if busy == nil {
		return nil, fmt.Errorf("invalid Busy pin name: %s", *busyPin)
	}

	log.Printf("Creating inky...")
	dev, err := inky.New(b, dc, reset, busy, opts)
	if err != nil {
		return nil, err
	}

	return dev, nil
}
