package files

import (
	"os"
	"io"
	"fmt"
	"github.com/Azat201003/summorist-mores/internal/config"
)

func ReadFile(moreId uint32, offset uint32, buffer []byte) (uint32, error) {
	conf := config.GetConfig()
	filePath := conf.FilePrefix + fmt.Sprintf("%d", moreId)  + conf.FilePostfix
	file, err := os.Open(filePath)
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
	conf := config.GetConfig()
	filePath := conf.FilePrefix + fmt.Sprintf("%d", moreId)  + conf.FilePostfix
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(buffer)
	return err
}

