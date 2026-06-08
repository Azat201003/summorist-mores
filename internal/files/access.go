package files

import (
	"fmt"
	"io"
	"os"
)

func getFilePath(moreId uint32) string {
	return os.Getenv("MORES_FILE_PREFIX") + fmt.Sprintf("%d", moreId) + os.Getenv("MORES_FILE_SUFFIX")
}

func ReadFile(moreId uint32, offset uint32, buffer []byte) (uint32, error) {
	file, err := os.Open(getFilePath(moreId))
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

func WriteFile(moreId uint32, buffer []byte) error {
	file, err := os.OpenFile(getFilePath(moreId), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(buffer)
	return err
}

func ClearFile(moreId uint32) error {
	file, err := os.OpenFile(getFilePath(moreId), os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}
