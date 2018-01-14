package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDailyStoryLine(t *testing.T) {
	sl, err := getDailyStoryline("access_token", "20151216")

	assert.Nil(t, err)
	assert.NotNil(t, sl)
	assert.True(t, false)
	fmt.Printf("%+v", sl)
}
