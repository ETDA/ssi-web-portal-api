package helpers

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestHashPassword(t *testing.T) {
	seedPassword := "testpassword"
	password, _ := HashPassword(seedPassword)
	assert.Equal(t, ComparePassword(*password, seedPassword), true)
	assert.Equal(t, ComparePassword(*password, seedPassword+"x"), false)
}
