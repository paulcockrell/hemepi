package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"

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

func generateImage(data *Response) (*image.Image, error) {
	fmt.Println("Generating image...")

	// Just put some text together for now
	var text = []string{}
	text = append(text, fmt.Sprintf("Metal: %s", data.Metal))
	text = append(text, fmt.Sprintf("Currency: %s", data.Currency))
	text = append(text, fmt.Sprintf("Low price: %.2f", data.LowPrice))
	text = append(text, fmt.Sprintf("High price: %.2f", data.HighPrice))

	// Read the font data
	fontBytes, err := ioutil.ReadFile(fontfile)
	if err != nil {
		return nil, err
	}

	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		pt.Y += c.PointToFixed(fontSize * spacing)
	}

	img := rgba.SubImage(rgba.Rect)

	return &img, nil
}
