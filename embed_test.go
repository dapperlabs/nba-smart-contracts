package nba

import "os"

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmbed(t *testing.T) {
	content, err := os.ReadFile("transactions/admin/fulfill_pack.cdc")
	assert.NoError(t, err)
	assert.Equal(t, AdminFulfillPack, content)
}
