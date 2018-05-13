package main

import "golang.org/x/tour/pic"

func Pic(dx, dy int) [][]uint8 {
	slice := make([][]uint8, dy)
	for row := range slice{
		slice[row] = make([]uint8, dx)
		for col := range slice[row]{
			slice[row][col] = uint8(row ^ col)
		}
	}
	return slice
}

func main() {
	pic.Show(Pic)
}
