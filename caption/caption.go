package caption // import "bygophers.com/go/memes/caption"

import (
	"image"
	"image/color"
	"image/draw"
)

// Caption describes the text overlay to draw on base image.
type Caption struct {
	Top    string
	Bottom string
}

// Draw a text overlay on a base image.
func Draw(base image.Image, req *Caption) (image.Image, error) {
	dst := image.NewRGBA(base.Bounds())
	draw.Draw(dst, dst.Bounds(), base, base.Bounds().Min, draw.Src)

	// TODO render text, instead of just ruining the image
	r := dst.Bounds().Inset(300)
	blue := color.RGBA{0, 0, 255, 255}
	draw.Draw(dst, r, &image.Uniform{blue}, image.ZP, draw.Src)

	return dst, nil
}
