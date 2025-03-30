package database

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewWithErrors(t *testing.T) {
	_, err := New(context.TODO(), nil, "")
	assert.Error(t, err)
}
