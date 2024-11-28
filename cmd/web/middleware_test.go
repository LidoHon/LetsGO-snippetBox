package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LidoHon/LetsGO-snippetBox.git/internal/assert"
)


func TestSecureHeaders(t *testing.T){
	rr :=httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)

	if err !=nil{
		t.Fatal(err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		w.Write([]byte("OK"))
	})
	secureHeaders(next).ServeHTTP(rr, r)
	rs := rr.Result()


	// Check that the middleware has correctly set the Content-Security-Policy
	// header on the response.
	expectedValue := "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com;"
	assert.Equal(t, rs.Header.Get("Content-Security-Policy"), expectedValue)

	// Check that the middleware has correctly set the Referrer-Policy
	// header on the response
	expectedValue ="origin-when-cross-origin"
	assert.Equal(t, rs.Header.Get("Referrer-Policy"), expectedValue)

	// Check that the middleware has correctly set the X-Frame-Options header
	// on the response.

	expectedValue = "nosniff"
	assert.Equal(t, rs.Header.Get("X-Content-Type-Options"), expectedValue)

	// Check that the middleware has correctly set the X-Frame-Options header
	// on the response.
	expectedValue = "deny"
	assert.Equal(t, rs.Header.Get("X-Frame-Options"), expectedValue)


	// Check that the middleware has correctly set the X-XSS-Protectionheader
	// on the response
	expectedValue ="0"
	assert.Equal(t, rs.Header.Get("X-XSS-Protection"), expectedValue)
	// Check that the middleware has correctly called the next handler inline
	// and the response status code and body are as expected.
	assert.Equal(t, rs.StatusCode, http.StatusOK)

	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	if err !=nil{
		t.Fatal(err)
	}
	bytes.TrimSpace(body)
	assert.Equal(t, string(body), "OK")

}

