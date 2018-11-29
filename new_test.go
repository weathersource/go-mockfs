package mockfs

import (
	"testing"

	assert "github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	client, server := New(t)
	assert.NotNil(t, client)
	assert.NotNil(t, server)
}
