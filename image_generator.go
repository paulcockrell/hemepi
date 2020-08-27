package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"

	"github.com/golang/freetype"
	"golang.org/x/image/font"
)

var (
	fontfile         = "./assets/luxisr.ttf"
	height           = 212
	width            = 104
	dpi      float64 = 72
	fontSize float64 = 12
	spacing          = 1.5

	black = color.RGBA{0, 0, 0, 0xff}
	white = color.RGBA{0xff, 0xff, 0xff, 0xff}
)

var text = []string{
	"Hello, Bob",
	"Hello, Moss",
	"Hello, Paul",
}

//func generateImage(data *Response, name string) error {
func generateImage() error {
	/*
		upLeft := image.Point{0, 0}
		lowRight := image.Point{width, height}
		img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	*/

	// Read the font data
	fontBytes, err := ioutil.ReadFile(fontfile)
	if err != nil {
		return err
	}

	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return err
	}

	// Initialize the context
	fg, bg := image.Black, image.White
	ruler := color.RGBA{0xdd, 0xdd, 0xdd, 0xff}
	rgba := image.NewRGBA(image.Rect(0, 0, 104, 212))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(f)
	c.SetFontSize(fontSize)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
	c.SetHinting(font.HintingNone)

	// Draw the guidelines
	for i := 0; i < 200; i++ {
		rgba.Set(10, 10+i, ruler)
		rgba.Set(10+i, 10, ruler)
	}

	// Draw the text
	pt := freetype.Pt(10, 10+int(c.PointToFixed(fontSize)>>6))
	for _, s := range text {
		_, err := c.DrawString(s, pt)
		if err != nil {
			return err
		}
		pt.Y += c.PointToFixed(fontSize * spacing)
	}

	// Save RGBA image to disk
	outFile, err := os.Create("out.png")
	if err != nil {
		return err
	}
	defer outFile.Close()

	b := bufio.NewWriter(outFile)
	err = png.Encode(b, rgba)
	if err != nil {
		return err
	}
	err = b.Flush()
	if err != nil {
		return err
	}

	fmt.Println("Image generated to file 'out.png'")
	return nil
}
