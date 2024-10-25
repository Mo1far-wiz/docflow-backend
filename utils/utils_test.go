package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "secret123"
	hashed, err := HashPassword(password)

	assert.NoError(t, err)
	assert.True(t, CheckPasswordHash(password, hashed))
}

func TestJWT(t *testing.T) {
	email := "john@example.com"
	userID := int64(1)

	token, err := GenerateToken(email, userID)
	assert.NoError(t, err)

	parsedID, err := VerifyToken(token)
	assert.NoError(t, err)
	assert.Equal(t, userID, parsedID)
}
