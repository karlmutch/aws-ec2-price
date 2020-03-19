package price

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrice(t *testing.T) {
	_, err := NewPricing("us-east-1")

	assert.Nil(t, err)
}
