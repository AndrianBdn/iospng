package main

import (
	"os"
	"image/png"
	"image/color"
	"fmt"
	"log"
	"github.com/nfnt/resize"
	"io"

	"bufio"
	"bytes"
	"encoding/binary"

	"hash/crc32"
	"compress/zlib"

	"io/ioutil"
)


func decodePngData(data []byte) []byte {

	var zdatbuf bytes.Buffer
	zdatbuf.Write([]byte{0x78, 0x1})
	zdatbuf.Write(data)
	zdatbuf.Write([]byte{0,0,0,0}) // don't know CRC, will get zlib.ErrChecksum

	reader, err := zlib.NewReader(&zdatbuf)
	if err != nil {
		panic(err)
	}

	dat, err := ioutil.ReadAll(reader)
	println(len(dat))

	if err != zlib.ErrChecksum {
		println("fail", err.Error())
		//panic(err)
	}
	reader.Close()
	return dat
	
}

func writeChunk(writer io.Writer, chunkType []byte, chunkData []byte, chunkCRC uint32, needCrc bool) {
	if (needCrc) {
		crc := crc32.NewIEEE()
		crc.Write(chunkType)
		crc.Write(chunkData)
		chunkCRC = crc.Sum32()
	}

	chunkLength := uint32(len(chunkData))
	binary.Write(writer, binary.BigEndian, &chunkLength)
	writer.Write(chunkType)
	writer.Write(chunkData)
	binary.Write(writer, binary.BigEndian, &chunkCRC)
}


func PngDecode(reader io.Reader, writer io.Writer) {
	bufReader := bufio.NewReader(reader)
	header := make([]byte, 8)


	bufReader.Read(header)

	if bytes.Compare([]byte("\x89PNG\r\n\x1a\n"), header) != 0 {
		panic("no png header")
	}

	writer.Write(header)

	// read chunk

	var w, h int;
	var datbuf bytes.Buffer

CHUNK:
	for {
		// read chunk
		var chunkLength uint32;
		binary.Read(bufReader, binary.BigEndian, &chunkLength)

		chunkType := make([]byte, 4)
		bufReader.Read(chunkType)
		//println(chunkLength)
		println(string(chunkType))

		chunkData := make([]byte, chunkLength)
		bufReader.Read(chunkData)

		var chunkCRC uint32;
		binary.Read(bufReader, binary.BigEndian, &chunkCRC)
		//



		switch {

		case bytes.Compare(chunkType, []byte("IHDR")) == 0:
			w = int(binary.BigEndian.Uint32(chunkData[:4]))
			h = int(binary.BigEndian.Uint32(chunkData[4:8]))
			println(w)
			println(h)


		case bytes.Compare(chunkType, []byte("CgBI")) == 0:
			continue;

		case bytes.Compare(chunkType, []byte("IDAT")) == 0:
			//tmp = true
			//if tmp {
				datbuf.Write(chunkData)
			//	println("buf ", datbuf.Len())
			//}
			//tmp = true





			//datbuf.Write(dat)

			continue;

		case bytes.Compare(chunkType, []byte("IEND")) == 0:

			raw := decodePngData(datbuf.Bytes())

			// size (w+1)*h
			for y := 0; y<h; y++ {

				for x := 0; x<w; x++   {

					row := y * w*4 + y;
					col := x*4 + 1

					b := raw[row + col + 0]
					r := raw[row + col + 2]

					raw[row + col + 0] = r
					raw[row + col + 2] = b

				}
			}


			var zdatbuf bytes.Buffer
			zwrite := zlib.NewWriter(&zdatbuf)
			zwrite.Write(raw)
			zwrite.Close()

			writeChunk(writer, []byte("IDAT"), zdatbuf.Bytes(), 0, true)
			writeChunk(writer, []byte("IEND"), []byte{}, 0, true)



			break CHUNK;
		}

		writeChunk(writer, chunkType, chunkData, chunkCRC, true)
	}
}


func main() {

	f, err := os.Open(os.Args[1])

	if err != nil {
		panic(err)
	}

	w, err := os.Create("./test.png")
	if err != nil {
		panic(err)
	}

	PngDecode(f, w);
	w.Close()


	f, err = os.Open("./test.png")

	if err != nil {
		panic(err)
	}


	bigImg, err := png.Decode(f)
	if err != nil {
		panic(err)
	}
	os.Exit(1)

	img := resize.Resize(80, 0, bigImg, resize.Lanczos3)


	if err != nil {
		log.Fatal(err)
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
