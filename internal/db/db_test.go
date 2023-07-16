package db

import (
	"testing"

	"github.com/edwintrumpet/prueba-tecnica-boletia/config"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	config.NewMock(1, 1, "")

	err := New()
	assert.Error(t, err)
}
