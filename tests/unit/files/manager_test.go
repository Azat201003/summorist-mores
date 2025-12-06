package files_test

import (
	"github.com/Azat201003/summorist-mores/internal/files"
	"github.com/stretchr/testify/assert"
	"testing"
	"math/rand/v2"
)

func TestWriteOk(t *testing.T) {
	str := []byte("Some testing stirng, that will be splitted")
	size := 4 // bytes per send
	count := (len(str)+size-1)/size
	moreId := rand.Uint32()
	for i := 0; i < count; i++ {
		err := files.WriteFile(moreId, str[i*size:min(len(str),(i+1)*size)])
		assert.NoError(t, err)
	}
}

