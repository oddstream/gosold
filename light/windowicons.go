package light

import (
	"bytes"
	_ "embed"
	"image"
	"log"
)

// https://www.iconsdb.com/white-icons/cards-icon.html

//go:embed icons/cards-16.png
var cards16IconBytes []byte

//go:embed icons/cards-32.png
var cards32IconBytes []byte

//go:embed icons/cards-48.png
var cards48IconBytes []byte

// WindowIcons create window icons in various resolutions.
// Exported so it can be used in package main.
func WindowIcons() []image.Image {
	var images []image.Image

	for _, b := range [][]byte{cards16IconBytes, cards32IconBytes, cards48IconBytes} {
		img, _, err := image.Decode(bytes.NewReader(b))
		if err != nil {
			log.Panic(err)
		}
		images = append(images, img)
	}

	return images
}
