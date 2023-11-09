package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pashagolub/pgxmock/v3"

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

func TestGetAllItems(t *testing.T) {
	// setup mock dependencies
	mockDBPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	deps := dependencies.NewDependencies(
		validator.New(),
		mockDBPool,
	)
	defer deps.CleanupDependencies()
	// define mock DB expectations
	rows := mockDBPool.NewRows([]string{"id", "uuid", "created_at", "name", "price"}).
		AddRow(
			1,
			"550e8400-e29b-41d4-a716-446655440000",
			time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
			"pi",
			float32(3.14),
		).
		AddRow(
			2,
			"550e8400-e29b-41d4-a716-446655440001",
			time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
			"tree-fiddy",
			float32(3.50),
		)
	mockDBPool.ExpectQuery("SELECT (.+) FROM item").
		WillReturnRows(rows)
	// setup router
	r := gin.Default()
	r.GET("/api/items/all", routes.HandleGetAllItems(deps))
	// exec request
	w := performRequest(r, "GET", "/api/items/all")
	// assert response code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}
	// assert full response body
	expected := `{"data":[{"id":1,"uuid":"550e8400-e29b-41d4-a716-446655440000","created_at":"2021-01-01T00:00:00Z","name":"pi","price":3.14},{"id":2,"uuid":"550e8400-e29b-41d4-a716-446655440001","created_at":"2021-01-01T00:00:00Z","name":"tree-fiddy","price":3.5}],"meta":{}}`
	if w.Body.String() != expected {
		t.Errorf("Expected %s, but got %s", expected, w.Body.String())
	}
	// assert db expectations were met
	if err := mockDBPool.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled DB expectations: %s", err)
	}
}

func TestGetItem(t *testing.T) {
	// setup mock dependencies
	mockDBPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	deps := dependencies.NewDependencies(
		validator.New(),
		mockDBPool,
	)
	defer deps.CleanupDependencies()
	// define mock DB expectations
	rows := mockDBPool.NewRows([]string{"id", "uuid", "created_at", "name", "price"}).
		AddRow(
			1,
			"550e8400-e29b-41d4-a716-446655440000",
			time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
			"pi",
			float32(3.14),
		)
	mockDBPool.ExpectQuery("SELECT (.+) FROM item WHERE id = (.+)").
		WithArgs(1).
		WillReturnRows(rows)
	// setup router
	r := gin.Default()
	r.GET("/api/items/:id", routes.HandleGetItem(deps))
	// exec request
	w := performRequest(r, "GET", "/api/items/1")
	// assert response code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}
	// assert full response body
	expected := `{"data":{"id":1,"uuid":"550e8400-e29b-41d4-a716-446655440000","created_at":"2021-01-01T00:00:00Z","name":"pi","price":3.14},"meta":{}}`
	if w.Body.String() != expected {
		t.Errorf("Expected %s, but got %s", expected, w.Body.String())
	}
	// assert db expectations were met
	if err := mockDBPool.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled DB expectations: %s", err)
	}
}
