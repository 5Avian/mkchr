package main

import (
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func writeTile(img image.Image, lookup map[color.Color]byte, tx int, ty int, writer io.Writer) error {
	result := make([]byte, 16)
	for y := 0; y < 8; y++ {
		result[y] = (lookup[img.At(tx+0, ty+y)]&0b01)<<7 |
			(lookup[img.At(tx+1, ty+y)]&0b01)<<6 |
			(lookup[img.At(tx+2, ty+y)]&0b01)<<5 |
			(lookup[img.At(tx+3, ty+y)]&0b01)<<4 |
			(lookup[img.At(tx+4, ty+y)]&0b01)<<3 |
			(lookup[img.At(tx+5, ty+y)]&0b01)<<2 |
			(lookup[img.At(tx+6, ty+y)]&0b01)<<1 |
			(lookup[img.At(tx+7, ty+y)]&0b01)<<0
	}
	for y := 0; y < 8; y++ {
		result[y+8] = (lookup[img.At(tx+0, ty+y)]&0b10)<<6 |
			(lookup[img.At(tx+1, ty+y)]&0b10)<<5 |
			(lookup[img.At(tx+2, ty+y)]&0b10)<<4 |
			(lookup[img.At(tx+3, ty+y)]&0b10)<<3 |
			(lookup[img.At(tx+4, ty+y)]&0b10)<<2 |
			(lookup[img.At(tx+5, ty+y)]&0b10)<<1 |
			(lookup[img.At(tx+6, ty+y)]&0b10)<<0 |
			(lookup[img.At(tx+7, ty+y)]&0b10)>>1
	}
	_, err := writer.Write(result)
	return err
}

func main() {
	if len(os.Args) != 2 {
		println("Usage: mkchr <file>")
		os.Exit(1)
	}
	inPath := os.Args[1]
	outPath := strings.TrimSuffix(inPath, filepath.Ext(inPath)) + ".chr"

	inFile, err := os.Open(inPath)
	if err != nil {
		println("input file error")
		println(err)
		os.Exit(1)
	}
	defer inFile.Close()

	img, _, err := image.Decode(inFile)
	if err != nil {
		println("image error")
		println(err)
		os.Exit(1)
	}

	bounds := img.Bounds()
	if bounds.Max.X%8 != 0 || bounds.Max.Y%8 != 0 {
		println("image resolution is not divisible by 8")
		os.Exit(1)
	}

	var nColors byte = 0
	lookup := map[color.Color]byte{}
	for y := 0; y < bounds.Max.Y; y++ {
		for x := 0; x < bounds.Max.X; x++ {
			color := img.At(x, y)
			if _, ok := lookup[color]; !ok {
				lookup[color] = nColors
				nColors++
			}
			if nColors > 4 {
				println("image contains more than 4 colors")
				os.Exit(1)
			}
		}
	}

	outFile, err := os.Create(outPath)
	if err != nil {
		println("output file error")
		println(err)
		os.Exit(1)
	}
	defer outFile.Close()

	for y := 0; y < bounds.Max.Y; y += 8 {
		for x := 0; x < bounds.Max.X; x += 8 {
			err = writeTile(img, lookup, x, y, outFile)
			if err != nil {
				println(err)
				os.Exit(1)
			}
		}
	}
}
