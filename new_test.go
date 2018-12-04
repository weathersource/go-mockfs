package mockfs

import (
	"testing"

	assert "github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	client, server, err := New()
	assert.NotNil(t, client)
	assert.NotNil(t, server)
	assert.Nil(t, err)
}
