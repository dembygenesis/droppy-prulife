package s3

import (
	"bytes"
	"errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"mime/multipart"
	"os"
)

func IsMultipartImage(f *multipart.FileHeader) error {
	file, err := f.Open()

	_, _, err = image.Decode(file)

	if err != nil {
		return errors.New("muiltipart is not an image file")
	}

	return nil
}

func IsImage(f *os.File)  {

}

func MultipartToReader(f *multipart.FileHeader) (*bytes.Reader, error) {
	multiPartFile, err := f.Open()

	if err != nil {
		return nil, err
	}

	size := f.Size
	buffer := make([]byte, size)
	multiPartFile.Read(buffer)

	conversion := bytes.NewReader(buffer)

	return conversion, err
}

func UploadToS3() {

}
