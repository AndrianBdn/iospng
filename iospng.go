package iospng

import (
	"io"
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"compress/zlib"
	"io/ioutil"
	"errors"
)

var (
	ErrPngHeader = errors.New("Not a Png");
	ErrImageData = errors.New("Unexpected amount of image data")
)


type pngChunk struct {
	chunkLength, chunkCRC uint32
	chunkType, chunkData []byte
}

func decodePngData(data []byte) ([]byte, error) {

	var zbuf bytes.Buffer
	zbuf.Write([]byte{0x78, 0x1}) 	// looks like a good zlib header
	zbuf.Write(data)
	zbuf.Write([]byte{0,0,0,0}) 	// don't know CRC, will get zlib.ErrChecksum

	reader, err := zlib.NewReader(&zbuf)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	dat, err := ioutil.ReadAll(reader)

	if err != zlib.ErrChecksum {
		return nil, err
	}
	return dat, nil
}

func (p *pngChunk) write(writer io.Writer, needCrc bool) error {
	if needCrc {
		crc := crc32.NewIEEE()
		crc.Write(p.chunkType)
		crc.Write(p.chunkData)
		p.chunkCRC = crc.Sum32()
	}

	chunkLength := uint32(len(p.chunkData))
	err := binary.Write(writer, binary.BigEndian, &chunkLength)
	if err != nil {
		return err
	}
	_, err = writer.Write(p.chunkType)
	if err != nil {
		return err
	}
	_, err = writer.Write(p.chunkData)
	if err != nil {
		return err
	}
	err = binary.Write(writer, binary.BigEndian, &p.chunkCRC)
	if err != nil {
		return err
	}
	return nil
}

func (p *pngChunk) read(reader io.Reader) error {

	if err := binary.Read(reader, binary.BigEndian, &p.chunkLength); err != nil {
		return err
	}

	p.chunkType = make([]byte, 4)

	if _, err := io.ReadFull(reader, p.chunkType); err != nil {
		return err
	}

	p.chunkData = make([]byte, p.chunkLength)

	if _, err := io.ReadFull(reader, p.chunkData); err != nil {
		return err
	}

	if err := binary.Read(reader, binary.BigEndian, &p.chunkCRC); err != nil {
		return err
	}

	return nil
}

func (p *pngChunk) is(kind string) bool {
	return string(p.chunkType) == kind
}

func rawImageFix(w, h int, raw []byte) error {
	if len(raw) != w*h*4 + h {
		return ErrImageData
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {

			row := y*w*4 + y;
			col := x*4 + 1

			b := raw[row+col+0]
			r := raw[row+col+2]

			raw[row+col+0] = r
			raw[row+col+2] = b

		}
	}
	return nil
}

func PngRevertOptimization(reader io.Reader, writer io.Writer) error {
	header := make([]byte, 8)
	if _, err := io.ReadFull(reader, header); err != nil {
		return err
	}

	if bytes.Compare([]byte("\x89PNG\r\n\x1a\n"), header) != 0 {
		return ErrPngHeader
	}

	writer.Write(header)

	var w, h int;
	var datbuf bytes.Buffer
	optimized := false

	for {
		var chunk pngChunk
		if err := chunk.read(reader); err != nil {
			return err
		}


		switch {

		case chunk.is("IHDR"):
			w = int(binary.BigEndian.Uint32(chunk.chunkData[:4]))
			h = int(binary.BigEndian.Uint32(chunk.chunkData[4:8]))


		case chunk.is("CgBI"):
			optimized = true
			continue;

		case chunk.is("IDAT"):
			if optimized {
				datbuf.Write(chunk.chunkData)
				continue;
			}


		case chunk.is("IEND"):
			if optimized {

				raw, err := decodePngData(datbuf.Bytes())
				if err != nil {
					return err
				}


				if err = rawImageFix(w, h, raw); err != nil {
					return err
				}

				var zdatbuf bytes.Buffer
				zwrite := zlib.NewWriter(&zdatbuf)
				zwrite.Write(raw)
				zwrite.Close()

				chunk.chunkType = []byte("IDAT")
				chunk.chunkData = zdatbuf.Bytes()
				err = chunk.write(writer, true)

				chunk.chunkType = []byte("IEND")
				chunk.chunkData = []byte{}
				err = chunk.write(writer, true)

				return nil
			}
		}

		if err := chunk.write(writer, false); err != nil {
			return nil
		}

	}

	return nil
}
