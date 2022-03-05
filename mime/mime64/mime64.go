package mime64

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
)

func Encode(content []byte) []byte {
	var mimeType = http.DetectContentType(content)
	var buf bytes.Buffer
	buf.WriteString("data:")
	buf.WriteString(mimeType)
	buf.WriteString(";base64,")
	buf.WriteString(base64.StdEncoding.EncodeToString(content))
	return buf.Bytes()
}

func EncodeToString(content []byte) string {
	return string(Encode(content))
}

func EncodeFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return EncodeToString(content), nil
}

func EncodeURL(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return EncodeToString(content), nil
}

func EncodeImagePNG(img image.Image) (string, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", err
	}
	return EncodeToString(buf.Bytes()), nil
}

func EncodeImageJPEG(img image.Image) (string, error) {
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, nil); err != nil {
		return "", err
	}
	return EncodeToString(buf.Bytes()), nil
}

func EncodeImageGIF(img image.Image) (string, error) {
	var buf bytes.Buffer
	if err := gif.Encode(&buf, img, nil); err != nil {
		return "", err
	}
	return EncodeToString(buf.Bytes()), nil
}
