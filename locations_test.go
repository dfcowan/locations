package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRound(t *testing.T) {
	assert.Equal(t, 0.12345, roundToFivePlaces(0.123454))
	assert.Equal(t, 0.12346, roundToFivePlaces(0.123456))

	assert.Equal(t, -0.12345, roundToFivePlaces(-0.123454))
	assert.Equal(t, -0.12346, roundToFivePlaces(-0.123456))
}
