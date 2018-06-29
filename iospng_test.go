package iospng

import (
	"bytes"
	"encoding/base64"
	"image/color"
	"image/png"
	"testing"
)

func TestGoodPng(t *testing.T) {
	// normal PNG file 3x3 pixels, does not contain CgBI
	const goodPngB64 = "iVBORw0KGgoAAAANSUhEUgAAAAMAAAADCAAAAABzQ+pjAAAAAXNSR0IArs4c6QAAAA" +
		"5JREFUCFtj/M/AwAjFABUPAwGHiY2AAAAAAElFTkSuQmCC"

	pngdata, _ := base64.StdEncoding.DecodeString(goodPngB64)
	reader := bytes.NewReader(pngdata)

	var w bytes.Buffer

	PngRevertOptimization(reader, &w)

	if bytes.Compare(pngdata, w.Bytes()) != 0 {
		t.Error("Good PNG was changed")
	}
}

func PerformRevertOptimizationTest(t *testing.T, base64data string, expectedBounds string) {
	pngdata, _ := base64.StdEncoding.DecodeString(base64data)
	reader := bytes.NewReader(pngdata)

	var w bytes.Buffer

	err := PngRevertOptimization(reader, &w)
	if err != nil {
		t.Error(err)
	}

	decReader := bytes.NewReader(w.Bytes())
	img, err := png.Decode(decReader) // crashes if PngRevertOptimization did wrong
	if err != nil {
		t.Error(err)
	}

	bString := img.Bounds().String()

	if bString != expectedBounds {
		t.Error("Expected ", bString, " to be ", expectedBounds)
	}

}

func TestIosPng(t *testing.T) {
	// 3x3 png optimized for iOS by Xcode
	const iOSPngB64 = "iVBORw0KGgoAAAAEQ2dCSVAAIAYsuHdmAAAADUlIRFIAAAADAAAAAwgGAAAAVii1vwAAAARnQU1BAACxjwv8YQUAA" +
		"AABc1JHQgCuzhzpAAAAIGNIUk0AAHomAACAhAAA+gAAAIDoAAB1MAAA6mAAADqYAAAXcJy6UTwAAAAJcEhZcwAACx" +
		"MAAAsTAQCanBgAAAAISURBVGP4jwQYcHIAIY6C+AAAAABJRU5ErkJggg=="

	PerformRevertOptimizationTest(t, iOSPngB64, "(0,0)-(3,3)")
}

func TestAlphaDemultiply(t *testing.T) {
	// 100x100 white png with 50 alpha
	const iOSPngB64 = "iVBORw0KGgoAAAAEQ2dCSVAAIAIr1bN/AAAADUlIRFIAAABkAAAAZAgGAAAAcOKVVAAAAARnQU1BAACxjwv8YQUA" +
		"AAABc1JHQgCuzhzpAAAAIGNIUk0AAHomAACAhAAA+gAAAIDoAAB1MAAA6mAAADqYAAAXcJy6UTwAAAAJcEhZcwAA" +
		"CxMAAAsTAQCanBgAAAAcaURPVAAAAAIAAAAAAAAAMgAAACgAAAAyAAAAMgAAAJD96hyjAAAAXElEQVTs0TEBAAAM" +
		"wjCkIx0LO3ekEppUr4oFQAQEiIAAERAgAgJEQAQEiIAAERAgAgJEQAQEiIAAERAgAgJEQAQEiIAAERAgAgJEQAQE" +
		"iIAAERAgAgJEQAQEiK4NAAD//yg9d44AAABYSURBVO3RMQEAAAzCMKQjHQs7d6QSmlSvigVABASIgAARECACAkRA" +
		"BASIgAARECACAkRABASIgAARECACAkRABASIgAARECACAkRABASIgAARECACAkRABASIrg2mAn5KAAAAAElFTkSuQmCC"

	pngdata, _ := base64.StdEncoding.DecodeString(iOSPngB64)
	reader := bytes.NewReader(pngdata)

	var w bytes.Buffer

	PngRevertOptimization(reader, &w)

	decReader := bytes.NewReader(w.Bytes())
	img, err := png.Decode(decReader)
	if err != nil {
		t.Error(err)
	}

	pixel := color.NRGBAModel.Convert(img.At(50, 50)).(color.NRGBA)

	if pixel.R != 0xff || pixel.G != 0xff || pixel.B != 0xff || pixel.A != 128 {
		t.Error("Expected color to be white@(128 alpha), got", pixel)
	}
}

func TestInterlacedPngDecode(t *testing.T) {
	const iOSPngB64 = "iVBORw0KGgoAAAAEQ2dCSVAAIAYsuHdmAAAADUlIRFIAAAAUAAAAFAgGAAAB+o4tmwAAAARnQU1BAACxjwv8YQUAAAABc1JHQgCuzhzpAAAAIGNIUk0AAHomAACAhAAA+gAAAIDoAAB1MAAA6mAAADqYAAAXcJy6UTwAAAQ7SURBVI2U3W8UZRSH35vtjX8DWSheQKIxLt8fkS5ose0aTbFCROveQCym7CqEbqFVEa4UBS1EiKghkgYjSALd0m1BabsUxZZKkfKh0tp2h1La7Xy8u7NfM48XYxeKkjjJk3cyc86Zec/v9x4hhMC2bXp7exGmaeJ2F+KPSEQ0GsUfkaiqipBSAiCWrI1R7J2DUFWV+qh0bnpvJ9nbeZSVF1vJZrNOvCegUN++E1+42okMhULYtg2AZVkUFBQgMpkMz6weoam9ht0nDSe9+w9JrP8oJSVl+JrfxjRNp6Y3EKXu/E6W77mFlP98veO2lv9TVVUR27Ztw7KsPFcG04S7EuRyuTyWZSHS6TTpdJr5byr093/Lvh2zefKJOfSfeoxVl07y04hJOp1GJBIJun6XzK+K8Vt/I+vXLMTtLuTSmVJKI7X4I5JkMonQNA1d12nolswLKngCCvM3D7G4dhR/RKLrOrquI1wuF5OTk0xOThKPx/+TlpYWREFBAb29vViWxdRl2za2bSOlJBQKIYRAPLjjZCrHiQsSI3V/t5ZlUVNT47Q8m81SVHuHFzbcpC36HqHOj/Fuv4aZyZLNZp0+ptNpth6ZwOcbYKDtNeY9PZeNr86irmUjZXv/Yqp9wjRNPAGFz1tPsGPTAtzuQoqLS+hqWYz3dBdjuulomEwm8QQUfHUtxK43MnPmbBYtmEv5mSDew31cG004fZRS8tZxFU+1QvmZLTR0fcSa5gAlrR9Q2SyR0kHouo6maZTsn8ATUJhXHcMTUKg4prGv22DqvZhyh6qqvBs18Eck/ojk/G09/7yzsxMRiUQeqcgULpfLkTAajWJZFrZt59f/w4M50WjUKdjT0zPNo7lcDpnMsvVInGXrYhR5R/C+OIC36gZF6//kpQ2jnL9h/isnl8vR09PjCJjJZMhkMugyTVHdKEvLY1Tt6ubm5U8ZiLzCrfBT3Giawc2muZxu9fPsz42s/OxXDvySyOdmMhlH6FQqRSqVQpcpFmy5w6LKGPWHzjFy9SvaDq9g1qzHcbsLqahYS0PDAWqrZvBjywpKzu7G++Vldl5M5F2TSqUcQySTSZouG3gCCosqFepPNqLc+p5vPinH7S6cxhsvz+bw2dfxhaspOuFMqOtjSUzTdIyTSCSQUnJPlSypHcUTUFi2bpCqtl30XXMKX72wn45T79DXd4hNZ7fjC1ez+oeDlH4xhj8iUeIPGcwwDAzD4OKQwaoP7+EJOCd7+aY+nj+2h7LwFnzhzZQ27+C5775m4fZBKo45U6lj8P6pNwwDoSgKmqblmVA19vforDuuUbx3nOXv32Vp3V1W7B6j7GCcymbHqPu6DcYntWm5iqIggsFgfpQ8jKqqDN9T6YtpXBnRGBpznP6o+GAwiBBC4HK5CAaDdHR0MDw8zMTEBPF4fNr6MOPj4wwNDdHe3k4gEHBMLQR/A97Pwk0AAAAASUVORK5CYII="

	PerformRevertOptimizationTest(t, iOSPngB64, "(0,0)-(20,20)")
}
