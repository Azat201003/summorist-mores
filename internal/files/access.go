// Package files is created for files manipulations
package files

import (
	"fmt"
	"io"
	"os"
)

func GetFilePath(moreId uint64) string {
	return os.Getenv("MORES_FILE_PREFIX") + fmt.Sprintf("%d", moreId) + os.Getenv("MORES_FILE_SUFFIX")
}

func ReadFile(moreId uint64, offset uint32, buffer []byte) (uint32, error) {
	file, err := os.Open(GetFilePath(moreId))
	if err != nil {
		return 0, err
	}
	defer file.Close()

	n, err := file.ReadAt(buffer, int64(offset))
	if err != nil && err != io.EOF {
		return uint32(n), err
	}

	return uint32(n), nil
}

func WriteFile(moreId uint64, buffer []byte) error {
	file, err := os.OpenFile(GetFilePath(moreId), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(buffer)
	return err
}

func ClearFile(moreId uint64) error {
	file, err := os.OpenFile(GetFilePath(moreId), os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

func RemoveFile(moreId uint64) error {
	return os.Remove(GetFilePath(moreId))
}
