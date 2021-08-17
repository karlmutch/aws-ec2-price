package price

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestPrice performing a trivial price retrival for
// the us-east-1 region
func TestPrice(t *testing.T) {
	_, err := NewPricing("us-east-1")

	assert.Nil(t, err)
}
