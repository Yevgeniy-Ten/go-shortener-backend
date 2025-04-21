package filestorage

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewFileStorage(t *testing.T) {
	tempFile, err := os.CreateTemp("", "test_storage_*.txt")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	fs, err := New(tempFile.Name(), nil)
	assert.NoError(t, err)
	assert.NotNil(t, fs)
	assert.NotNil(t, fs.file)
	assert.NotNil(t, fs.writer)
	assert.NotNil(t, fs.scanner)
	assert.NotNil(t, fs.storage)

	fs, err = New("", nil)
	assert.Error(t, err)
	assert.Nil(t, fs)
}
