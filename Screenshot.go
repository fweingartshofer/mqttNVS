package main

import (
	"github.com/vova616/screenshot"
	"image"
	"image/png"
	"os"
)

func main() {
	captureScreen()
}

func captureDisplay() {
	img, err := screenshot.CaptureRect(image.Rect(0, 0, 1920, 1080))
	if err != nil {
		panic(err)
	}
	f, err := os.Create("./ss.png")
	if err != nil {
		panic(err)
	}
	err = png.Encode(f, img)
	if err != nil {
		panic(err)
	}
	f.Close()
}
