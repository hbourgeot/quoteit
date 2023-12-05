package generators

import (
	"bytes"
	"fmt"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"strings"
)

func wrapText(text string, maxWidth int, face font.Face) []string {
	var lines []string
	var line string

	words := strings.Fields(text)
	spaceWidth := font.MeasureString(face, " ").Round()

	for _, word := range words {
		// Medir el ancho del word actual
		lineWidth := font.MeasureString(face, line+word).Round()

		if len(line)+len(word)+1 > maxWidth && line != "" {
			lines = append(lines, line)
			line = word // Comienza una nueva línea
		} else {
			if line != "" {
				line += " "
				lineWidth += spaceWidth
			}
			line += word
		}
	}

	if line != "" {
		lines = append(lines, line)
	}

	return lines
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

func LoadFont(variant string) (font.Face, *truetype.Font, error) {
	var file string
	switch variant {
	case "Black", "BlackItalic", "Bold", "BoldItalic", "Italic", "Light", "LightItalic", "Medium", "MediumItalic", "Regular", "Thin", "ThinItalic":
		file = fmt.Sprintf("assets/fonts/Roboto-%s.ttf", variant)
	default:
		log.Println("No font provided or found")
		return nil, nil, fmt.Errorf("No font provided or found")
	}
	fontBytes, err := os.ReadFile(file)
	if err != nil {
		return nil, nil, err
	}

	fontParsed, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, nil, err
	}

	fontFace := truetype.NewFace(fontParsed, &truetype.Options{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingNone,
	})

	return fontFace, fontParsed, nil
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
	_, fnt, err := LoadFont("Regular") // Assuming LoadFont can handle "Regular"
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
	// Cargar las fuentes
	_, boldFontTtf, err := LoadFont("Bold")
	if err != nil {
		return nil, err
	}
	regularFont, regularFontTtf, err := LoadFont("Regular")
	if err != nil {
		return nil, err
	}

	if profileImg == nil {
		// Crea y dibuja la imagen de las iniciales si no hay imagen de perfil
		profileImg = roundProfileImage(createInitialsImage(name))
	} else {
		profileImg = roundProfileImage(profileImg)
	}

	profileWidth := profileImg.Bounds().Dx()
	profileRect := image.Rect(10, 10, 10+profileWidth, 190) // Padding de 10 píxeles
	nameWidth, _ := calculateTextWidth(name, boldFontTtf, 24)
	messageWidth, _ := calculateTextWidth(message, regularFontTtf, 24)

	wrappedMessage := wrapText(message, 50, regularFont) // Dibujar la imagen de perfil
	wrappedLen := len(wrappedMessage)

	var totalWidth int
	if len(message) > 40 {
		totalWidth = profileRect.Max.X + 790
	} else {
		totalWidth = max(profileRect.Max.X+10+nameWidth, profileRect.Max.X+10+messageWidth) + 50
	}

	totalHeight := ((3 + wrappedLen) * 27) + 20

	rgba := image.NewRGBA(image.Rect(0, 0, totalWidth, max(totalHeight, 200)))

	draw.Draw(rgba, profileRect, profileImg, image.Point{}, draw.Over)

	textBox := image.Rect(profileRect.Max.X+10, 10, totalWidth-10, max(totalHeight, 200)) // Dibujar el fondo con bordes redondeados
	y := textBox.Min.Y + 24
	drawRoundedRectangle(rgba, textBox, color.Black, 20)
	// Dibujar el texto del nombre
	drawText(rgba, boldFontTtf, name, image.Pt(textBox.Min.X+10, y), color.White) // Ajusta el color y la posición según sea necesario

	y += 40
	for _, line := range wrappedMessage {
		drawText(rgba, regularFontTtf, line, image.Pt(textBox.Min.X+10, y), color.White)
		y += 26
	}

	// Encode the image as PNG
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, rgba); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
