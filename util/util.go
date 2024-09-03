package util

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

const (
	serverStorageDir = "./server/storage/"
	clientStorageDir = "./client/storage/"
)

func SplitFileIntoChunks(fileName string, chunkSize int) ([][]byte, error) {
	file, err := os.Open(serverStorageDir + fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var chunks [][]byte
	buf := make([]byte, chunkSize)

	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if n == 0 {
			break
		}
		chunk := make([]byte, n)
		copy(chunk, buf[:n])
		chunks = append(chunks, chunk)
	}

	return chunks, nil
}

func ReadDirectory() ([]os.DirEntry, error) {
	files, err := os.ReadDir(serverStorageDir)
	if err != nil {
		logrus.Error("Error reading directory: -> ", err)
		return nil, err
	}
	return files, nil
}
func SaveToStorage(err error, chunks [][]byte, fileName string) {

	fo, err := os.Create(clientStorageDir + fileName)
	defer fo.Close()
	if err != nil {
		logrus.Error("Error creating file: -> ", err)
	}
	for _, chunk := range chunks {
		_, err = fo.Write(chunk)
		if err != nil {
			logrus.Error("Error writing file: -> ", err)
			return
		}
	}
	defer fo.Close()

}
