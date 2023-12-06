package generators

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
)

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
	fontSize := 30.0

	var totalWidth int
	if len(message) > 40 {
		fontSize = 20
		totalWidth = profileRect.Max.X + 550
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
	drawText(rgba, boldFontTtf, name, image.Pt(textBox.Min.X+10, y), color.White, fontSize) // Ajusta el color y la posición según sea necesario

	y += 40
	for _, line := range wrappedMessage {
		drawText(rgba, regularFontTtf, line, image.Pt(textBox.Min.X+10, y), color.White, fontSize)
		y += 26
	}

	// Encode the image as PNG
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, rgba); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
