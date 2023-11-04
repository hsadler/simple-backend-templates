package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"example-server/routes"
)

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestStatus(t *testing.T) {
	r := gin.Default()
	r.GET("/status", routes.HandleStatus)

	w := performRequest(r, "GET", "/status")

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}

	expected := `{"status":"ok"}`
	if w.Body.String() != expected {
		t.Errorf("Expected %s, but got %s", expected, w.Body.String())
	}
}

func TestMetrics(t *testing.T) {
	r := gin.Default()
	r.GET("/metrics", routes.HandleMetrics(r))

	w := performRequest(r, "GET", "/metrics")

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}

	expectedSubstring := "go_"
	if !strings.Contains(w.Body.String(), expectedSubstring) {
		t.Errorf("Expected response body to contain substring %s, but got %s", expectedSubstring, w.Body.String())
	}
}
