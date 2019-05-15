package utils_test

import (
	"github.com/illuminati1911/goira/internal/utils"
	"github.com/illuminati1911/goira/testutils"
	"testing"
)

func TestMin(t *testing.T) {
	assert := testutils.NewAssert(t)
	assert.Equals(utils.Min(1, 3), 1)
	assert.Equals(utils.Min(-500, -2), -500)
	assert.Equals(utils.Min(1, 99999), 1)
	assert.Equals(utils.Min(-1, 0), -1)
	assert.Equals(utils.Min(0, -1), -1)
	assert.Equals(utils.Min(0, 0), 0)
	assert.Equals(utils.Min(999999, 999999), 999999)
}

func TestMax(t *testing.T) {
	assert := testutils.NewAssert(t)
	assert.Equals(utils.Max(-1, 5), 5)
	assert.Equals(utils.Max(-1, -5), -1)
	assert.Equals(utils.Max(0, 0), 0)
	assert.Equals(utils.Max(5, 1000), 1000)
	assert.Equals(utils.Max(-1000, 1000), 1000)
	assert.Equals(utils.Max(-1, 5), 5)
}

func TestClamp(t *testing.T) {
	assert := testutils.NewAssert(t)
	assert.Equals(utils.Clamp(0, 3, 1), 1)
	assert.Equals(utils.Clamp(0, 3, 10), 3)
	assert.Equals(utils.Clamp(0, 3, -1), 0)
	assert.Equals(utils.Clamp(-500, 45, -200), -200)
	assert.Equals(utils.Clamp(-500, 45, 400), 45)
	assert.Equals(utils.Clamp(-500, 45, 0), 0)
	assert.Equals(utils.Clamp(-500, 45, -25500), -500)
}

func TestReverse(t *testing.T) {
	assert := testutils.NewAssert(t)
	assert.Equals(utils.Reverse(0x55), byte(0xAA))
	assert.Equals(utils.Reverse(0x1), byte(0x80))
	assert.Equals(utils.Reverse(0xF), byte(0xF0))
	assert.Equals(utils.Reverse(0x2), byte(0x40))
	assert.Equals(utils.Reverse(0xFF), byte(0xFF))
	assert.Equals(utils.Reverse(0x0), byte(0x0))
}
