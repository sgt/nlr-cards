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

func TestNlr_FindLastIdInASmartWay(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	nlr := NewNLR()
	lastCardNumber, err := nlr.FindLastCardNumberInASmartWay(42)
	assert.NoError(t, err)
	assert.Equal(t, 39, lastCardNumber)
}
