package fn

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidatePassword(t *testing.T) {
	type testCases struct {
		testCase string
		pass     bool
	}

	values := []testCases{
		{"$adigunHAMMED001", true},
		{"adigunHAMMED001", false},
		{"$adigunHAMMED", false},
		{"$HAMMED001", false},
		{"$aHAM1", false},
	}
	for _, next := range values {
		err := ValidatePassword(next.testCase)
		assert.Equal(t, next.pass, err == nil)
	}
}
