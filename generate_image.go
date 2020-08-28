package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/disintegration/imaging"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

var (
	fontfile         = "./assets/Verdana.ttf"
	width            = 212
	height           = 104
	upLeft           = image.Point{0, 0}
	lowRight         = image.Point{width, height}
	dpi      float64 = 72
	fontSize float64 = 10
	spacing          = 1.25

	black = color.RGBA{0, 0, 0, 0xff}
	white = color.RGBA{0xff, 0xff, 0xff, 0xff}
)

func generateImage(data *Response) (*image.Image, error) {
	log.Println("Generating image...")

	font, err := loadFont("Verdana.ttf")
	if err != nil {
		return nil, err
	}

	baseImage, err := loadImage("baseImage.png")
	if err != nil {
		return nil, err
	}

	// Initialize image
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(rgba, rgba.Bounds(), *baseImage, image.Point{0, 0}, draw.Over)

	// Setup font context
	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFontSize(fontSize)
	c.SetFont(font)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(image.Black)

	// Convert raw data object to display text
	lines := BuildLines(data)
	if err != nil {
		return nil, err
	}

	for _, line := range lines {
		c.SetFontSize(line.fontSize)
		pt := freetype.Pt(line.x, line.y)
		_, err = c.DrawString(line.text, pt)
		if err != nil {
			return nil, err
		}
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

func assetPath() (string, error) {
	// Get executable directory
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	dir := path.Dir(ex)

	return fmt.Sprintf("%s/assets/", dir), nil
}

func loadFont(fontFile string) (*truetype.Font, error) {
	ap, err := assetPath()
	if err != nil {
		return nil, err
	}

	// Read the font data
	fontBytes, err := ioutil.ReadFile(ap + fontFile)
	if err != nil {
		return nil, err
	}

	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}

	return font, nil
}

func loadImage(imageFile string) (*image.Image, error) {
	ap, err := assetPath()
	if err != nil {
		return nil, err
	}

	baseImageFile, err := os.Open(ap + imageFile)
	if err != nil {
		return nil, err
	}
	defer baseImageFile.Close()

	baseImage, err := png.Decode(baseImageFile)
	if err != nil {
		return nil, err
	}

	return &baseImage, nil
}

// Line holds information for rendering text to a display
type Line struct {
	x, y     int
	text     string
	fontSize float64
}

// BuildLines converts a Response into an array of lines
func BuildLines(data *Response) []Line {
	// Metal/Currency pair title
	title := Line{
		fontSize: 22,
		x:        109,
		y:        27,
		text:     fmt.Sprintf("%s%s", data.Metal, data.Currency),
	}

	ounces := Line{
		fontSize: 18,
		x:        10,
		y:        52,
		text:     fmt.Sprintf("%.2f", data.Price),
	}

	kilos := Line{
		fontSize: 18,
		x:        10,
		y:        76,
		text:     fmt.Sprintf("%.2f", data.Price*32.15),
	}

	// Price change from previous close price
	difference := data.Price - data.PrevClosePrice
	percentage := (difference / data.PrevClosePrice) * 100
	sign := "+"
	if data.Price < data.PrevClosePrice {
		sign = "-"
	}
	percentChange := Line{
		fontSize: 18,
		x:        107,
		y:        52,
		text:     fmt.Sprintf("%s%.2f%%", sign, percentage),
	}

	priceChange := Line{
		fontSize: 18,
		x:        114,
		y:        73,
		text:     fmt.Sprintf("[%.2f]", difference),
	}

	timeT := time.Unix(data.Timestamp, 0)
	updateTime := Line{
		fontSize: 9,
		x:        112,
		y:        93,
		text:     timeT.Format("02/01/06 15:04:05"),
	}

	return []Line{
		title,
		ounces,
		kilos,
		percentChange,
		priceChange,
		updateTime,
	}
}
