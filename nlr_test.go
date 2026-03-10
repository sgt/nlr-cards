package nlr_cards

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNlr_saveDir(t *testing.T) {
	nlr := NewNLR()
	assert.Equal(t, "downloads/0/0", nlr.saveDir(1))
	assert.Equal(t, "downloads/1/2", nlr.saveDir(1234))
}
