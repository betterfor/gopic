package resize

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
)

type CompressType string

const (
	CompressNone   CompressType = ""
	CompressAuto                = "auto"
	CompressCustom              = "custom"
)

func Compress(t CompressType, base int, file string) (io.Reader, error) {
	size, img, err := getReadSizeFile(file)
	if err != nil {
		return nil, err
	}
	var x, y int
	switch t {
	case CompressAuto:
		if size < base {
			base = 1
		} else {
			base = size / base
		}
		x, y = img.Bounds().Dx()/base, img.Bounds().Dy()/base
	case CompressCustom:
		x, y = img.Bounds().Dx()/base, img.Bounds().Dy()/base
	default:
		x, y = img.Bounds().Dx(), img.Bounds().Dy()
	}
	fmt.Println(base)

	var buf = bytes.NewBuffer(nil)
	img = resize.Thumbnail(uint(x), uint(y), img, resize.Lanczos3)
	err = jpeg.Encode(buf, img, nil)
	return buf, err
}

func getReadSizeFile(file string) (size int, img image.Image, err error) {
	fh, err := os.Open(file)
	if err != nil {
		return
	}
	defer fh.Close()

	// get file size
	fileInfo, err := fh.Stat()
	if err != nil {
		return
	}

	switch filepath.Ext(file) {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(fh)
	case ".png":
		img, err = png.Decode(fh)
	case ".gif":
		img, err = gif.Decode(fh)
	default:
		err = errors.New("no support picture extension")
	}
	return int(fileInfo.Size()), img, err
}
