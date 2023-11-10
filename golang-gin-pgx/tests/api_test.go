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
	"example-server/models"
	"example-server/routes"
)

// MOCKS

const (
	mockRecord1 = "mockRecord1"
	mockRecord2 = "mockRecord2"
)

var mockRecords = map[string]models.Item{
	mockRecord1: {
		ID:        1,
		UUID:      "550e8400-e29b-41d4-a716-446655440000",
		CreatedAt: time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
		Name:      "pi",
		Price:     float32(3.14),
	},
	mockRecord2: {
		ID:        2,
		UUID:      "550e8400-e29b-41d4-a716-446655440001",
		CreatedAt: time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
		Name:      "tree-fiddy",
		Price:     float32(3.50),
	},
}

// HELPERS

func getMockDependencies() (*dependencies.Dependencies, pgxmock.PgxPoolIface) {
	// setup mock dependencies
	mockDBPool, err := pgxmock.NewPool()
	if err != nil {
		panic(err)
	}
	deps := dependencies.NewDependencies(
		validator.New(),
		mockDBPool,
	)
	return deps, mockDBPool
}

func getMockRows(mockDBPool pgxmock.PgxPoolIface, items []models.Item) *pgxmock.Rows {
	// define mock DB expectations
	rows := mockDBPool.NewRows([]string{"id", "uuid", "created_at", "name", "price"})
	for _, item := range items {
		rows.AddRow(
			item.ID,
			item.UUID,
			item.CreatedAt,
			item.Name,
			item.Price,
		)
	}
	return rows
}

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// TESTS

func TestStatus(t *testing.T) {
	r := gin.Default()
	r.GET("/status", routes.HandleStatus)
	w := performRequest(r, "GET", "/status")
	expectedStatusCode := http.StatusOK
	if w.Code != expectedStatusCode {
		t.Errorf("Expected status code %d, but got %d", expectedStatusCode, w.Code)
	}
	expectedBody := `{"status":"ok"}`
	if w.Body.String() != expectedBody {
		t.Errorf("Expected %s, but got %s", expectedBody, w.Body.String())
	}
}

func TestMetrics(t *testing.T) {
	r := gin.Default()
	r.GET("/metrics", routes.HandleMetrics(r))
	w := performRequest(r, "GET", "/metrics")
	expectedStatusCode := http.StatusOK
	if w.Code != expectedStatusCode {
		t.Errorf("Expected status code %d, but got %d", expectedStatusCode, w.Code)
	}
	expectedSubstring := "go_"
	if !strings.Contains(w.Body.String(), expectedSubstring) {
		t.Errorf("Expected response body to contain substring %s, but got %s", expectedSubstring, w.Body.String())
	}
}

func TestGetAllItems200(t *testing.T) {
	// setup mock dependencies
	deps, mockDBPool := getMockDependencies()
	rows := getMockRows(mockDBPool, []models.Item{mockRecords[mockRecord1], mockRecords[mockRecord2]})
	mockDBPool.ExpectQuery("SELECT (.+) FROM item").
		WillReturnRows(rows)
	// setup router
	r := gin.Default()
	r.GET("/api/items/all", routes.HandleGetAllItems(deps))
	// exec request
	w := performRequest(r, "GET", "/api/items/all")
	// assert response code
	expectedStatusCode := http.StatusOK
	if w.Code != expectedStatusCode {
		t.Errorf("Expected status code %d, but got %d", expectedStatusCode, w.Code)
	}
	// assert full response body
	expectedBody := `{"data":[{"id":1,"uuid":"550e8400-e29b-41d4-a716-446655440000","created_at":"2021-01-01T00:00:00Z","name":"pi","price":3.14},{"id":2,"uuid":"550e8400-e29b-41d4-a716-446655440001","created_at":"2021-01-01T00:00:00Z","name":"tree-fiddy","price":3.5}],"meta":{}}`
	if w.Body.String() != expectedBody {
		t.Errorf("Expected %s, but got %s", expectedBody, w.Body.String())
	}
	// assert db expectations were met
	if err := mockDBPool.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled DB expectations: %s", err)
	}
}

func TestGetItem200(t *testing.T) {
	// setup mock dependencies
	deps, mockDBPool := getMockDependencies()
	rows := getMockRows(mockDBPool, []models.Item{mockRecords[mockRecord1]})
	mockDBPool.ExpectQuery("SELECT (.+) FROM item WHERE id = (.+)").
		WithArgs(1).
		WillReturnRows(rows)
	// setup router
	r := gin.Default()
	r.GET("/api/items/:id", routes.HandleGetItem(deps))
	// exec request
	w := performRequest(r, "GET", "/api/items/1")
	// assert response code
	expectedStatusCode := http.StatusOK
	if w.Code != expectedStatusCode {
		t.Errorf("Expected status code %d, but got %d", expectedStatusCode, w.Code)
	}
	// assert full response body
	expectedBody := `{"data":{"id":1,"uuid":"550e8400-e29b-41d4-a716-446655440000","created_at":"2021-01-01T00:00:00Z","name":"pi","price":3.14},"meta":{}}`
	if w.Body.String() != expectedBody {
		t.Errorf("Expected %s, but got %s", expectedBody, w.Body.String())
	}
	// assert db expectations were met
	if err := mockDBPool.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled DB expectations: %s", err)
	}
}

func TestGetItem404(t *testing.T) {
	// setup mock dependencies
	deps, mockDBPool := getMockDependencies()
	rows := getMockRows(mockDBPool, []models.Item{})
	mockDBPool.ExpectQuery("SELECT (.+) FROM item WHERE id = (.+)").
		WithArgs(1).
		WillReturnRows(rows)
	// setup router
	r := gin.Default()
	r.GET("/api/items/:id", routes.HandleGetItem(deps))
	// exec request
	w := performRequest(r, "GET", "/api/items/1")
	// assert response code
	expectedStatusCode := http.StatusNotFound
	if w.Code != expectedStatusCode {
		t.Errorf("Expected status code %d, but got %d", expectedStatusCode, w.Code)
	}
	// assert full response body
	expectedBody := `{"error":"Item not found"}`
	if w.Body.String() != expectedBody {
		t.Errorf("Expected %s, but got %s", expectedBody, w.Body.String())
	}
	// assert db expectations were met
	if err := mockDBPool.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled DB expectations: %s", err)
	}
}

func TestGetItem400(t *testing.T) {
	// setup mock dependencies
	deps, _ := getMockDependencies()
	// setup router
	r := gin.Default()
	r.GET("/api/items/:id", routes.HandleGetItem(deps))
	// exec request
	w := performRequest(r, "GET", "/api/items/invalid")
	// assert response code
	expectedStatusCode := http.StatusBadRequest
	if w.Code != expectedStatusCode {
		t.Errorf("Expected status code %d, but got %d", expectedStatusCode, w.Code)
	}
	// assert full response body
	expectedBody := `{"error":"Invalid Item ID"}`
	if w.Body.String() != expectedBody {
		t.Errorf("Expected %s, but got %s", expectedBody, w.Body.String())
	}
}

func TestGetItems200(t *testing.T) {
	// setup mock dependencies
	deps, mockDBPool := getMockDependencies()
	rows := getMockRows(mockDBPool, []models.Item{mockRecords[mockRecord1], mockRecords[mockRecord2]})
	mockDBPool.ExpectQuery("SELECT (.+) FROM item WHERE id = ANY(.+)").
		WithArgs([]int{1, 2}).
		WillReturnRows(rows)
	// setup router
	r := gin.Default()
	r.GET("/api/items", routes.HandleGetItems(deps))
	// exec request
	w := performRequest(r, "GET", "/api/items?item_ids=1&item_ids=2")
	// assert response code
	expectedStatusCode := http.StatusOK
	if w.Code != expectedStatusCode {
		t.Errorf("Expected status code %d, but got %d", expectedStatusCode, w.Code)
	}
	// assert full response body
	expectedBody := `{"data":[{"id":1,"uuid":"550e8400-e29b-41d4-a716-446655440000","created_at":"2021-01-01T00:00:00Z","name":"pi","price":3.14},{"id":2,"uuid":"550e8400-e29b-41d4-a716-446655440001","created_at":"2021-01-01T00:00:00Z","name":"tree-fiddy","price":3.5}],"meta":{}}`
	if w.Body.String() != expectedBody {
		t.Errorf("Expected %s, but got %s", expectedBody, w.Body.String())
	}
	// assert db expectations were met
	if err := mockDBPool.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled DB expectations: %s", err)
	}
}

func TestGetItems400MissingItemIds(t *testing.T) {
	// setup mock dependencies
	deps, _ := getMockDependencies()
	// setup router
	r := gin.Default()
	r.GET("/api/items", routes.HandleGetItems(deps))
	// exec request
	w := performRequest(r, "GET", "/api/items")
	// assert response code
	expectedStatusCode := http.StatusBadRequest
	if w.Code != expectedStatusCode {
		t.Errorf("Expected status code %d, but got %d", expectedStatusCode, w.Code)
	}
	// assert full response body
	expectedBody := `{"error":"Missing Item IDs"}`
	if w.Body.String() != expectedBody {
		t.Errorf("Expected %s, but got %s", expectedBody, w.Body.String())
	}
}

func TestGetItems400InvalidItemIds(t *testing.T) {
	// setup mock dependencies
	deps, _ := getMockDependencies()
	// setup router
	r := gin.Default()
	r.GET("/api/items", routes.HandleGetItems(deps))
	// exec request
	w := performRequest(r, "GET", "/api/items?item_ids=invalid")
	// assert response code
	expectedStatusCode := http.StatusBadRequest
	if w.Code != expectedStatusCode {
		t.Errorf("Expected status code %d, but got %d", expectedStatusCode, w.Code)
	}
	// assert full response body
	expectedBody := `{"error":"Invalid Item ID"}`
	if w.Body.String() != expectedBody {
		t.Errorf("Expected %s, but got %s", expectedBody, w.Body.String())
	}
}
