package useful

import (
	"bufio"
	"compress/gzip"
	"io/ioutil"
	"mime/multipart"
	"os"
)

func CompressAndSave(incomingFile *multipart.FileHeader, path string) error {
	f, err := incomingFile.Open()

	if err != nil {
		return err
	}

	read := bufio.NewReader(f)

	data, err := ioutil.ReadAll(read)

	if err != nil {
		return err
	}

	file, err := os.Create(path)

	if err != nil {
		return err
	}

	w := gzip.NewWriter(file)

	w.Write(data)

	w.Close()

	return nil
}
