package auth

import (
	"testing"

	"github.com/google/uuid"
)

const testSecret = "secret-key"

var testToken string

func init() {
	userID := uuid.New().String()
	testToken, _ = GenerateTestToken(testSecret, userID, "user")
}

func BenchmarkJWTValidation(b *testing.B) {

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ValidateTestToken(testToken, testSecret)
	}

}

func TestValidateToken(t *testing.T) {
	_, err := ValidateTestToken(testToken, testSecret)
	if err != nil {
		t.Fatalf("Token Validation failed unexpectedly: %v", err)
	}
}
