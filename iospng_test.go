package iospng

import (
	"testing"
	"encoding/base64"
	"bytes"
	"image/png"
)


func TestGoodPng(t *testing.T) {
	// normal PNG file 3x3 pixels, does not contain CgBI
	const goodPngB64 = "iVBORw0KGgoAAAANSUhEUgAAAAMAAAADCAAAAABzQ+pjAAAAAXNSR0IArs4c6QAAAA" +
			   "5JREFUCFtj/M/AwAjFABUPAwGHiY2AAAAAAElFTkSuQmCC"

	pngdata, _ := base64.StdEncoding.DecodeString(goodPngB64)
	reader := bytes.NewReader(pngdata)

	var w bytes.Buffer;

	PngRevertOptimization(reader, &w)

	if bytes.Compare(pngdata, w.Bytes()) != 0 {
		t.Error("Good PNG was changed")
	}
}

func TestIosPng(t *testing.T) {
	// 3x3 png optimized for iOS by Xcode
	const iOSPngB64 = "iVBORw0KGgoAAAAEQ2dCSVAAIAYsuHdmAAAADUlIRFIAAAADAAAAAwgGAAAAVii1vwAAAARnQU1BAACxjwv8YQUAA" +
		          "AABc1JHQgCuzhzpAAAAIGNIUk0AAHomAACAhAAA+gAAAIDoAAB1MAAA6mAAADqYAAAXcJy6UTwAAAAJcEhZcwAACx" +
		  	  "MAAAsTAQCanBgAAAAISURBVGP4jwQYcHIAIY6C+AAAAABJRU5ErkJggg=="

	pngdata, _ := base64.StdEncoding.DecodeString(iOSPngB64)
	reader := bytes.NewReader(pngdata)

	var w bytes.Buffer;

	PngRevertOptimization(reader, &w)

	decReader := bytes.NewReader(w.Bytes())
	img, err := png.Decode(decReader) // crashes if PngRevertOptimization did wrong
	if err != nil {
		t.Error(err)
	}

	bString := img.Bounds().String()
	exp := "(0,0)-(3,3)"

	if bString != exp {
		t.Error("Expected ", bString, " to be ", exp)
	}
}