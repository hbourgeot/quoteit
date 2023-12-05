package generators

import (
	"fmt"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
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
