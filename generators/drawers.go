package generators

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"image"
	"image/color"
	"image/draw"
	"log"
)

func drawText(img *image.RGBA, fnt *truetype.Font, text string, position image.Point, clr color.Color, fontSize float64) {
	// drawText draws the text on img using the provided font
	// Set up freetype context with the font and draw the text
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(fnt)
	c.SetFontSize(fontSize)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.NewUniform(clr))
	pt := freetype.Pt(position.X, position.Y+int(c.PointToFixed(fontSize)>>6))
	_, err := c.DrawString(text, pt)
	if err != nil {
		log.Println(err)
	}
}

func roundProfileImage(src image.Image) *image.RGBA {
	// El radio será la mitad del ancho o del alto, lo que sea menor
	size := src.Bounds().Size()
	r := min(size.X, size.Y) / 2
	mask := image.NewRGBA(src.Bounds())

	// Pintar la máscara de transparente
	draw.Draw(mask, mask.Bounds(), image.Transparent, image.Point{}, draw.Src)

	// Dibujar un círculo blanco en la máscara
	for y := -r; y < r; y++ {
		for x := -r; x < r; x++ {
			if x*x+y*y <= r*r {
				mask.Set(size.X/2+x, size.Y/2+y, color.Opaque)
			}
		}
	}

	// Crear la imagen final con el perfil redondeado
	dst := image.NewRGBA(src.Bounds())
	draw.DrawMask(dst, dst.Bounds(), src, image.Point{}, mask, image.Point{}, draw.Over)
	return dst
}

func createInitialsImage(name string) *image.RGBA {
	fontSize := 24.0
	bgColor := color.RGBA{70, 130, 180, 255}
	textColor := color.White

	// Create an image
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	// Load the font
	_, fnt, err := LoadFont("Regular") // Assuming LoadFont can handle "Regular"
	if err != nil {
		log.Println("Failed to load font:", err)
		return img
	}

	// Get the initials
	initials := getInitials(name)

	// Draw the initials in the center of the image
	pt := image.Pt(100/2-len(initials)*int(fontSize)/4, 100/2+int(fontSize)/2) // Centering logic
	drawText(img, fnt, initials, pt, textColor, fontSize)

	return img
}

func drawRoundedRectangle(dst *image.RGBA, r image.Rectangle, clr color.Color, radius int) {
	// Rellena el centro
	fill := image.NewUniform(clr)
	draw.Draw(dst, image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Max.Y), fill, image.Point{}, draw.Src)

	// Rellena los lados rectos
	draw.Draw(dst, image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Min.Y), fill, image.Point{}, draw.Src)
	draw.Draw(dst, image.Rect(r.Min.X, r.Max.Y, r.Max.X, r.Max.Y), fill, image.Point{}, draw.Src)
	draw.Draw(dst, image.Rect(r.Min.X, r.Min.Y, r.Min.X, r.Max.Y), fill, image.Point{}, draw.Src)
	draw.Draw(dst, image.Rect(r.Max.X, r.Min.Y, r.Max.X, r.Max.Y), fill, image.Point{}, draw.Src)
}
