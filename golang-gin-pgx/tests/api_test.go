package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"example-server/database"
	"example-server/dependencies"
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

func TestGetItem(t *testing.T) {
	r := gin.Default()
	// WORKING
	mockDBPool, err := database.SetupTestDB()
	if err != nil {
		t.Fatal(err)
	}
	deps := dependencies.NewDependencies(
		validator.New(),
		mockDBPool,
	)
	defer deps.CleanupDependencies()
	r.GET("/api/items/:id", routes.HandleGetItem(deps))

	w := performRequest(r, "GET", "/api/items/1")

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}

	expectedSubstring := `"id":1`
	if !strings.Contains(w.Body.String(), expectedSubstring) {
		t.Errorf("Expected response body to contain substring %s, but got %s", expectedSubstring, w.Body.String())
	}
}
