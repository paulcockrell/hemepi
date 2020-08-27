package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"log"

	"github.com/disintegration/imaging"
	"github.com/golang/freetype"
	"golang.org/x/image/font"
)

var (
	fontfile         = "./assets/luxisr.ttf"
	height           = 104
	width            = 212
	dpi      float64 = 72
	fontSize float64 = 14
	spacing          = 1.25

	black = color.RGBA{0, 0, 0, 0xff}
	white = color.RGBA{0xff, 0xff, 0xff, 0xff}
)

func generateImage(data *Response) (*image.Image, error) {
	log.Println("Generating image...")

	// Just put some text together for now
	var text = []string{}
	text = append(text, fmt.Sprintf("Metal: %s, Currency: %s", data.Metal, data.Currency))
	text = append(text, fmt.Sprintf("Low: %.2f, High: %.2f", data.LowPrice, data.HighPrice))
	text = append(text, fmt.Sprintf("Ask: %.2f, Bid: %.2f", data.Ask, data.Bid))

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
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))
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

	// If you look at the PI with external ports at bottom, then the display
	// will be presented to you with the short side on the Y axis, and the
	// data cable on the right.. So we need to rotate 90 counter clockwise
	// and horizontally flip the image so it appears as one would expect
	rotImg := imaging.Rotate90(rgba)
	mirrorImg := imaging.FlipH(rotImg)
	subImg := mirrorImg.SubImage(mirrorImg.Rect)

	return &subImg, nil
}
