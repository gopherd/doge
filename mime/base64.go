package mime

import (
	"bytes"
	"encoding/base64"
	"io"
	"net/http"
	"os"
)

func EncodeBase64(content []byte) []byte {
	var mimeType = http.DetectContentType(content)
	var buf bytes.Buffer
	buf.WriteString("data:")
	buf.WriteString(mimeType)
	buf.WriteString(";base64,")
	buf.WriteString(base64.StdEncoding.EncodeToString(content))
	return buf.Bytes()
}

func EncodeBase64ToString(content []byte) string {
	return string(EncodeBase64(content))
}

func EncodeBase64File(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return EncodeBase64ToString(content), nil
}

func EncodeBase64URL(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return EncodeBase64ToString(content), nil
}
