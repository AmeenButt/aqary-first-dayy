package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAuthMiddleware(t *testing.T) {
	os.Setenv("JWT_SECRET", "ksdjvjksdn::DAvcdsw3r3rdsfdsfvQ@!^%&%&^*&^%&fsaesvsddcvsvzadsnrdfbfdzbfdbfzvb.dfbfdbfd.dvsfbsdfbfdsbdfzxfbd")

	// Test case: Valid token
	validToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOjF9.t-HAl7Gk94QTjMUQu_jP9ukoT8FHDryxxDFXxXbWYYY"
	mockContextValid, _ := createMockContext(validToken)

	// Call AuthMiddleware directly with the mock context
	AuthMiddleware()(mockContextValid)
	t.Log(mockContextValid.Writer.Status())
	// Check the response
	if mockContextValid.Writer.Status() != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, mockContextValid.Writer.Status())
	}

	// Test case: Invalid token
	invalidToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOjF9.t-HAl7Gk94QTjMUQu_jP9ukoT8FHDryxxDFXxXbWYYd"
	mockContextInvalid, _ := createMockContext(invalidToken)

	// Call AuthMiddleware directly with the mock context
	AuthMiddleware()(mockContextInvalid)

	// Check the response
	if mockContextInvalid.Writer.Status() != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, mockContextInvalid.Writer.Status())
	}
}

func createMockContext(authorizationHeader string) (*gin.Context, *httptest.ResponseRecorder) {
	// Create a mock Gin context
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", authorizationHeader)
	context, _ := gin.CreateTestContext(w)
	context.Request = req

	return context, w
}
