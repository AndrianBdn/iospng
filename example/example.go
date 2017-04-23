package main

import (
	"image/png"
	"image/color"
	"fmt"
	"bytes"
	"github.com/andrianbdn/iospng"
	"os"
)


func main() {

	src, err := os.Open("./test.png")
	defer src.Close()

	var w bytes.Buffer
	iospng.PngRevertOptimization(src, &w)

	pngReader := bytes.NewReader(w.Bytes())

	img, err := png.Decode(pngReader)
	if err != nil {
		panic(err)
	}

	levels := []string{" ", "░", "▒", "▓", "█"}

	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			c := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
			level := c.Y / 51 // 51 * 5 = 255
			if level == 5 {
				level--
			}
			fmt.Print(levels[level])
		}
		fmt.Print("\n")
	}
}
