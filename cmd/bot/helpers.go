package main

import (
	"bytes"
	"fmt"
	"github.com/hbourgeot/quoteme/tgbot"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
)

func getImageFormat(r io.Reader) (string, error) {
	buf := make([]byte, 512) // 512 bytes should be enough for the magic number
	if _, err := r.Read(buf); err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buf)
	switch contentType {
	case "image/jpeg":
		return "jpeg", nil
	case "image/png":
		return "png", nil
	// add other cases as needed
	default:
		return "", fmt.Errorf("unrecognized image format")
	}
}

func decodeImage(r io.Reader, format string) (image.Image, error) {
	switch format {
	case "jpeg":
		return jpeg.Decode(r)
	case "png":
		return png.Decode(r)
	// add other cases as needed
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

func getImage(bot *tgbot.BotAPI, fileConf tgbot.FileConfig) image.Image {
	file, err := bot.GetFile(fileConf)
	if err != nil {
		log.Fatal("linea 18", err)
	}

	url, err := bot.GetFileDirectURL(file.FileID)
	if err != nil {
		log.Fatal("linea 23", err)
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("linea 28", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("bad status: %s", resp.Status)
	}

	// Read the entire response body
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error leyendo")
		return nil
	}

	// Determine the image format
	format, err := getImageFormat(bytes.NewReader(data))
	if err != nil {
		log.Println("Error obteniendo el formato")
		return nil
	}

	// Decode the image
	img, err := decodeImage(bytes.NewReader(data), format)
	if err != nil {
		log.Println("Error decodificando")
		return nil
	}

	return img
}
