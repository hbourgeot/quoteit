package generators

import (
	"bytes"
	"fmt"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"strings"
)

func LoadFont(variant string) (*truetype.Font, error) {
	var file string
	switch variant {
	case "Black", "BlackItalic", "Bold", "BoldItalic", "Italic", "Light", "LightItalic", "Medium", "MediumItalic", "Regular", "Thin", "ThinItalic":
		file = fmt.Sprintf("assets/fonts/Roboto-%s.ttf", variant)
	default:
		log.Println("No font provided or found")
		return nil, fmt.Errorf("No font provided or found")
	}
	fontBytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	fontParsed, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}

	return fontParsed, nil
}

func drawText(img *image.RGBA, fnt *truetype.Font, text string, position image.Point, clr color.Color) {
	fontSize := 30.0
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

func createInitialsImage(name string) *image.RGBA {
	fontSize := 24.0
	bgColor := color.RGBA{70, 130, 180, 255}
	textColor := color.White

	// Create an image
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	// Load the font
	fnt, err := LoadFont("Regular") // Assuming LoadFont can handle "Regular"
	if err != nil {
		log.Println("Failed to load font:", err)
		return img
	}

	// Get the initials
	initials := getInitials(name)

	// Draw the initials in the center of the image
	pt := image.Pt(100/2-len(initials)*int(fontSize)/4, 100/2+int(fontSize)/2) // Centering logic
	drawText(img, fnt, initials, pt, textColor)

	return img
}

func drawRoundedRectangle(img *image.RGBA, rect image.Rectangle, clr color.Color, radius int) {
	// Create a mask image to draw rounded corners
	mask := image.NewRGBA(img.Bounds())
	draw.Draw(mask, mask.Bounds(), &image.Uniform{color.Transparent}, image.Point{}, draw.Src)

	// Define the corner arc
	quarterCircle := func(r int) []image.Point {
		points := []image.Point{}
		for y := -r; y <= r; y++ {
			for x := -r; x <= r; x++ {
				if x*x+y*y <= r*r {
					points = append(points, image.Point{x, y})
				}
			}
		}
		return points
	}(radius)

	// Draw four corners
	corners := []image.Point{
		rect.Min,
		image.Pt(rect.Max.X-radius, rect.Min.Y),
		image.Pt(rect.Min.X, rect.Max.Y-radius),
		rect.Max.Sub(image.Pt(radius, radius)),
	}

	for _, corner := range corners {
		for _, p := range quarterCircle {
			mask.Set(corner.X+p.X, corner.Y+p.Y, clr)
		}
	}

	// Fill the centers
	centerRects := []image.Rectangle{
		image.Rect(rect.Min.X+radius, rect.Min.Y, rect.Max.X-radius, rect.Min.Y+radius), // Top
		image.Rect(rect.Min.X, rect.Min.Y+radius, rect.Max.X, rect.Max.Y-radius),        // Center
		image.Rect(rect.Min.X+radius, rect.Max.Y-radius, rect.Max.X-radius, rect.Max.Y), // Bottom
	}
	for _, r := range centerRects {
		draw.Draw(mask, r, &image.Uniform{clr}, image.Point{}, draw.Src)
	}

	// Draw the mask onto the original image
	draw.Draw(img, rect, mask, rect.Min, draw.Over)
}

func getInitials(name string) string {
	parts := strings.Fields(name)
	initials := ""
	for i := 0; i < len(parts) && i < 2; i++ {
		if len(parts[i]) > 0 {
			initials += strings.ToUpper(parts[i][:1])
		}
	}
	return initials
}

func calculateTextWidth(text string, fnt *truetype.Font, fontSize float64) (int, error) {
	// Configurar un contexto de freetype para medir el texto
	c := freetype.NewContext()
	c.SetFont(fnt)
	c.SetFontSize(fontSize)

	// Medir el ancho del texto
	opts := truetype.Options{Size: fontSize}
	face := truetype.NewFace(fnt, &opts)

	width := 0
	for _, x := range text {
		awidth, ok := face.GlyphAdvance(rune(x))
		if !ok {
			return 0, fmt.Errorf("could not calculate advance width for rune %v", x)
		}
		width += awidth.Ceil()
	}

	return width, nil
}

func GenerateImage(profileImg image.Image, name, message string) ([]byte, error) {
	// Create a new image with a transparent background
	rgba := image.NewRGBA(image.Rect(0, 0, 800, 200)) // Adjust dimensions as needed

	var profileRect image.Rectangle
	// Draw the profile image or initials image
	// Asegúrate de que profileImg no sea nil
	if profileImg != nil {
		// Suponiendo que quieres la imagen de perfil en la esquina superior izquierda
		profileRect = image.Rect(0, 0, profileImg.Bounds().Dx(), profileImg.Bounds().Dy())
		draw.Draw(rgba, profileRect, profileImg, image.Point{}, draw.Over)
	} else {
		// Crea y dibuja la imagen de las iniciales si no hay imagen de perfil
		profileImg := createInitialsImage(name)
		profileRect = image.Rect(0, 0, profileImg.Bounds().Dx(), profileImg.Bounds().Dy())
		draw.Draw(rgba, profileRect, profileImg, image.Point{}, draw.Over)
	}

	boldFont, err := LoadFont("Bold")
	if err != nil {
		return nil, err
	}

	regularFont, err := LoadFont("Regular")
	if err != nil {
		return nil, err
	}

	messageWidth, err := calculateTextWidth(message, regularFont, 30)

	// Ajustar las posiciones de los recuadros basándose en las anchuras calculadas
	nameBox := image.Rect(profileRect.Max.X+10, 5, profileRect.Max.X+10+messageWidth+20, 80)                   // +20 para el padding
	messageBox := image.Rect(nameBox.Min.X, nameBox.Max.Y-10, nameBox.Min.X+messageWidth+20, nameBox.Max.Y+80) // +10 para espacio entre recuadros

	// Dibujar recuadro con bordes redondeados para el nombre
	drawRoundedRectangle(rgba, nameBox, color.Black, 0) // Ajusta el color y radio según sea necesario
	// Dibujar el texto del nombre
	drawText(rgba, boldFont, name, nameBox.Min.Add(image.Pt(10, 20)), color.White) // Ajusta el color y la posición según sea necesario

	// Dibujar recuadro con bordes redondeados para el mensaje
	drawRoundedRectangle(rgba, messageBox, color.Black, 0)
	// Dibujar el mensaje
	drawText(rgba, regularFont, message, messageBox.Min.Add(image.Pt(10, 20)), color.White) // Ajusta el color y la posición según sea necesario

	// Encode the image as PNG
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, rgba); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
