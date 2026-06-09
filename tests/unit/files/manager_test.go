package files_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/Azat201003/summorist-mores/internal/files"
	"github.com/stretchr/testify/assert"
	// "math/rand/v2"
)

func TestWriteReadOk(t *testing.T) {
	// Prepare
	os.Setenv("MORES_FILE_PREFIX", "../testdata/files_")
	os.Setenv("MORES_FILE_SUFFIX", ".txt")

	str := []byte("Some testing stirng, that will be splitted")
	size := uint32(4) // bytes per send
	ln := uint32(len(str))
	count := (ln + size - 1) / size
	moreId := uint64(1)

	assert.NoError(t, files.ClearFile(moreId))

	// Write tests
	var i uint32
	for i = 0; i < count; i++ {
		err := files.WriteFile(moreId, str[i*size:min(ln, (i+1)*size)])
		assert.NoError(t, err, "Some error, while appending part of file")
	}

	filepath := os.Getenv("MORES_FILE_PREFIX") + fmt.Sprintf("%d", moreId) + os.Getenv("MORES_FILE_SUFFIX")
	assert.FileExists(t, filepath)
	data, err := os.ReadFile(filepath)
	assert.NoError(t, err, "Cannot read file")
	assert.Equal(t, str, data, "File contains wrong data")

	// Read tests
	var result []byte
	for i = 0; i < count-uint32(1); i++ {
		buffer := make([]byte, size)
		n, err := files.ReadFile(moreId, i*size, buffer)
		result = append(result, buffer...)
		assert.NoError(t, err, "Error reading file part")
		assert.Equal(t, n, size, "Read bytes count mismatch for full chunk")
	}
	buffer := make([]byte, ln%size)
	n, err := files.ReadFile(moreId, i*size, buffer)
	assert.NoError(t, err, "Error reading last file part")
	assert.Equal(t, n, ln%size, "Read bytes count mismatch for last chunk")
	result = append(result, buffer...)

	assert.Equal(t, str, result, "Assembled data does not match original")
}
