package utils

import (
	"crypto/sha256"
	"fmt"
	"os"
)

func GetFileData(filePath string) (fileData []byte, err error) {
	fileData, err = os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return fileData, nil
}

func GetFileHash(filePath string) (fileHash string, err error) {
	var (
		fileData []byte
	)

	fileData, err = GetFileData(filePath)
	if err != nil {
		return fileHash, err
	}
	fileHash, err = GetHash(fileData)
	Must(err)

	return fileHash, err
}

func GetHash(fileData []byte) (fileHash string, err error) {
	sum := sha256.Sum256(fileData)
	return fmt.Sprintf("%x", sum), err
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}
