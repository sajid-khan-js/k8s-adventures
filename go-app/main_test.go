package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNamespacesRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/namespaces", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	// Crude, but kube-system is always there so reasonable smoke test
	assert.Contains(t, w.Body.String(), "kube-system")
}

// TODO write more tests for 100% coverage and use mocks
