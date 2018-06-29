package main

import (
	"bytes"
	"fmt"
	"github.com/andrianbdn/iospng"
	"image/png"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("usage: %s <source.png> <dest.png>\n", os.Args[0])
		return
	}

	srcPath := os.Args[1]
	dstPath := os.Args[2]

	src, err := os.Open(srcPath)
	if err != nil {
		log.Fatalln("cannot open source file: ", err.Error())
	}
	defer src.Close()

	if _, err := os.Stat(dstPath); err == nil {
		log.Fatalln("destination file exists, delete it first")
	}

	var w bytes.Buffer
	err = iospng.PngRevertOptimization(src, &w)
	if err != nil {
		log.Fatalln("PngRevertOptimization error", err.Error())
	}

	pngReader := bytes.NewReader(w.Bytes())

	_, err = png.Decode(pngReader)
	if err != nil {
		log.Fatalln("png decode check error", err.Error())
	}

	dst, err := os.Create(dstPath)
	if err != nil {
		log.Fatalln("cannot open destination file: ", err.Error())
	}
	defer src.Close()

	_, err = dst.Write(w.Bytes())

	if err != nil {
		log.Fatalln("cannot write to destination file: ", err.Error())
	}

	log.Println("conversion done, no errors")
}
