package file

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"regexp"
)

func GetFileContentType(out *bytes.Reader) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

func GetMultiPartAsBuffer(f *multipart.FileHeader) []byte {
	multipartFile, _ := f.Open()

	size := f.Size
	buffer := make([]byte, size)
	multipartFile.Read(buffer)

	return buffer
}

func GetMultiPartFileType(f *multipart.FileHeader) string {

	fileType := http.DetectContentType(GetMultiPartAsBuffer(f))

	reDb := regexp.MustCompile("image/")
	trimmedFileType := reDb.ReplaceAllString(fileType, "")

	return trimmedFileType
}

func Dump(i interface{}) {
	res2B, _ := json.Marshal(i)
	fmt.Println(string(res2B))
}