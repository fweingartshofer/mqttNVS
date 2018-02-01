package main

import (
	"image"
	"os"
	"image/png"
	"github.com/vova616/screenshot"
)

func main(){
	img, err := screenshot.CaptureRect(image.Rect(0,0,1920, 1080))
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
