package main

import (
	"golang.org/x/tour/pic"
	"image/color"
	"image"
)

type Image struct{
	width int
	height int
	colorA uint8
	colorB uint8
}

func (img Image) Bounds() image.Rectangle{
	return image.Rect(0,0,img.width,img.height)
}

func (img Image) ColorModel() color.Model{
	return color.RGBAModel
}


func (img Image) At(x, y int) color.Color{
	return color.RGBA{uint8(x), uint8(y), img.colorA, img.colorB}
}

func main() {
	m := Image{100,100,253,254}
	pic.ShowImage(m)
}
