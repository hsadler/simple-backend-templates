package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
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

func performRequest(r http.Handler, method string, path string, body ...string) *httptest.ResponseRecorder {
	var req *http.Request
	if len(body) > 0 {
		req, _ = http.NewRequest(method, path, strings.NewReader(body[0]))
	} else {
		req, _ = http.NewRequest(method, path, nil)
	}
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
	// setup mock dependencies and DB query expectations
	deps, mockDBPool := getMockDependencies()
	rows := getMockRows(mockDBPool, []models.Item{mockRecords[mockRecord1], mockRecords[mockRecord2]})
	mockDBPool.ExpectQuery("SELECT (.+) FROM item ORDER BY id OFFSET (.+) LIMIT (.)").
		WithArgs(0, 2).
		WillReturnRows(rows)
	// setup router
	r := gin.Default()
	r.GET("/api/items/all", routes.HandleGetAllItems(deps))
	// exec request
	w := performRequest(r, "GET", "/api/items/all?offset=0&chunkSize=2")
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

func TestGetAllItems200Empty(t *testing.T) {
	// setup mock dependencies and DB query expectations
	deps, mockDBPool := getMockDependencies()
	rows := getMockRows(mockDBPool, nil)
	mockDBPool.ExpectQuery("SELECT (.+) FROM item ORDER BY id OFFSET (.+) LIMIT (.)").
		WithArgs(0, 2).
		WillReturnRows(rows)
	// setup router
	r := gin.Default()
	r.GET("/api/items/all", routes.HandleGetAllItems(deps))
	// exec request
	w := performRequest(r, "GET", "/api/items/all?offset=0&chunkSize=2")
	// assert response code
	expectedStatusCode := http.StatusOK
	if w.Code != expectedStatusCode {
		t.Errorf("Expected status code %d, but got %d", expectedStatusCode, w.Code)
	}
	// assert full response body
	expectedBody := `{"data":[],"meta":{}}`
	if w.Body.String() != expectedBody {
		t.Errorf("Expected %s, but got %s", expectedBody, w.Body.String())
	}
	// assert db expectations were met
	if err := mockDBPool.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled DB expectations: %s", err)
	}
}

func TestGetAllItems400InvalidChunkSize(t *testing.T) {
	// setup mock dependencies
	deps, _ := getMockDependencies()
	// setup router
	r := gin.Default()
	r.GET("/api/items/all", routes.HandleGetAllItems(deps))
	// exec request
	w := performRequest(r, "GET", "/api/items/all?offset=5&chunkSize=0")
	// assert response code
	expectedStatusCode := http.StatusBadRequest
	if w.Code != expectedStatusCode {
		t.Errorf("Expected status code %d, but got %d", expectedStatusCode, w.Code)
	}
}

func TestGetAllItems400MissingQueryParams(t *testing.T) {
	// setup mock dependencies
	deps, _ := getMockDependencies()
	// setup router
	r := gin.Default()
	r.GET("/api/items/all", routes.HandleGetAllItems(deps))
	// exec request
	w := performRequest(r, "GET", "/api/items/all")
	// assert response code
	expectedStatusCode := http.StatusBadRequest
	if w.Code != expectedStatusCode {
		t.Errorf("Expected status code %d, but got %d", expectedStatusCode, w.Code)
	}
}

func TestGetAllItems500PostgresError(t *testing.T) {
	// setup mock dependencies and DB query expectations
	deps, mockDBPool := getMockDependencies()
	mockDBPool.ExpectQuery("SELECT (.+) FROM item ORDER BY id OFFSET (.+) LIMIT (.)").
		WithArgs(0, 2).
		WillReturnError(&pgconn.PgError{Code: "12345"})
	// setup router
	r := gin.Default()
	r.GET("/api/items/all", routes.HandleGetAllItems(deps))
	// exec request
	w := performRequest(r, "GET", "/api/items/all?offset=0&chunkSize=2")
	// assert response code
	expectedStatusCode := http.StatusInternalServerError
	if w.Code != expectedStatusCode {
		t.Errorf("Expected status code %d, but got %d", expectedStatusCode, w.Code)
	}
}

func TestGetItem200(t *testing.T) {
	// setup mock dependencies and DB query expectations
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
	// setup mock dependencies and DB query expectations
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

func TestGetItem500PostgresError(t *testing.T) {
	// setup mock dependencies and DB query expectations
	deps, mockDBPool := getMockDependencies()
	mockDBPool.ExpectQuery("SELECT (.+) FROM item WHERE id = (.+)").
		WithArgs(1).
		WillReturnError(&pgconn.PgError{Code: "12345"})
	// setup router
	r := gin.Default()
	r.GET("/api/items/:id", routes.HandleGetItem(deps))
	// exec request
	w := performRequest(r, "GET", "/api/items/1")
	// assert response code
	expectedStatusCode := http.StatusInternalServerError
	if w.Code != expectedStatusCode {
		t.Errorf("Expected status code %d, but got %d", expectedStatusCode, w.Code)
	}
}

func TestGetItems200(t *testing.T) {
	// setup mock dependencies and DB query expectations
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

func TestGetItems500PostgresError(t *testing.T) {
	// setup mock dependencies and DB query expectations
	deps, mockDBPool := getMockDependencies()
	mockDBPool.ExpectQuery("SELECT (.+) FROM item WHERE id = ANY(.+)").
		WithArgs([]int{1, 2}).
		WillReturnError(&pgconn.PgError{Code: "12345"})
	// setup router
	r := gin.Default()
	r.GET("/api/items", routes.HandleGetItems(deps))
	// exec request
	w := performRequest(r, "GET", "/api/items?item_ids=1&item_ids=2")
	// assert response code
	expectedStatusCode := http.StatusInternalServerError
	if w.Code != expectedStatusCode {
		t.Errorf("Expected status code %d, but got %d", expectedStatusCode, w.Code)
	}
}

func TestCreateItem201(t *testing.T) {
	// setup mock dependencies and DB query expectations
	deps, mockDBPool := getMockDependencies()
	mockCreateRecord := mockRecords[mockRecord1]
	rows := getMockRows(mockDBPool, []models.Item{mockCreateRecord})
	mockDBPool.ExpectQuery("INSERT INTO item (.+) VALUES (.+) RETURNING id").
		WithArgs(mockCreateRecord.Name, mockCreateRecord.Price).
		WillReturnRows(mockDBPool.NewRows([]string{"id"}).AddRow(mockCreateRecord.ID))
	mockDBPool.ExpectQuery("SELECT (.+) FROM item WHERE id = (.+)").
		WithArgs(mockCreateRecord.ID).
		WillReturnRows(rows)
	// setup router
	r := gin.Default()
	r.POST("/api/items", routes.HandleCreateItem(deps))
	// exec request
	createItemRequest := models.CreateItemRequest{
		Data: models.ItemIn{
			Name:  mockCreateRecord.Name,
			Price: mockCreateRecord.Price,
		},
	}
	createItemRequestJson, _ := json.Marshal(createItemRequest)
	w := performRequest(r, "POST", "/api/items", string(createItemRequestJson))
	// assert response code
	expectedStatusCode := http.StatusCreated
	if w.Code != expectedStatusCode {
		t.Errorf("Expected status code %d, but got %d", expectedStatusCode, w.Code)
	}
	// assert full response body
	expectedBody := `{"data":{"id":1,"uuid":"550e8400-e29b-41d4-a716-446655440000","created_at":"2021-01-01T00:00:00Z","name":"pi","price":3.14},"meta":{"created":true}}`
	if w.Body.String() != expectedBody {
		t.Errorf("Expected %s, but got %s", expectedBody, w.Body.String())
	}
	// assert db expectations were met
	if err := mockDBPool.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled DB expectations: %s", err)
	}
}

func TestCreateItem400InvalidJson(t *testing.T) {
	// setup mock dependencies
	deps, _ := getMockDependencies()
	// setup router
	r := gin.Default()
	r.POST("/api/items", routes.HandleCreateItem(deps))
	// exec request
	w := performRequest(r, "POST", "/api/items", `{"data":{"invalid":"payload"}}`)
	// assert response code
	expectedStatusCode := http.StatusBadRequest
	if w.Code != expectedStatusCode {
		t.Errorf("Expected status code %d, but got %d", expectedStatusCode, w.Code)
	}
}

func TestCreateItem400InvalidItemIn(t *testing.T) {
	// setup mock dependencies
	deps, _ := getMockDependencies()
	// setup router
	r := gin.Default()
	r.POST("/api/items", routes.HandleCreateItem(deps))
	// exec request
	createItemRequest := models.CreateItemRequest{
		Data: models.ItemIn{
			Name:  "invalid price",
			Price: float32(-1),
		},
	}
	createItemRequestJson, _ := json.Marshal(createItemRequest)
	w := performRequest(r, "POST", "/api/items", string(createItemRequestJson))
	// assert response code
	expectedStatusCode := http.StatusBadRequest
	if w.Code != expectedStatusCode {
		t.Errorf("Expected status code %d, but got %d", expectedStatusCode, w.Code)
	}
}

func TestCreateItem409Duplicate(t *testing.T) {
	// setup mock dependencies and DB query expectations
	deps, mockDBPool := getMockDependencies()
	mockCreateRecord := mockRecords[mockRecord1]
	mockDBPool.ExpectQuery("INSERT INTO item (.+) VALUES (.+) RETURNING id").
		WithArgs(mockCreateRecord.Name, mockCreateRecord.Price).
		WillReturnError(&pgconn.PgError{Code: "23505"})
	// setup router
	r := gin.Default()
	r.POST("/api/items", routes.HandleCreateItem(deps))
	// exec request
	createItemRequest := models.CreateItemRequest{
		Data: models.ItemIn{
			Name:  mockCreateRecord.Name,
			Price: mockCreateRecord.Price,
		},
	}
	createItemRequestJson, _ := json.Marshal(createItemRequest)
	w := performRequest(r, "POST", "/api/items", string(createItemRequestJson))
	// assert response code
	expectedStatusCode := http.StatusConflict
	if w.Code != expectedStatusCode {
		t.Errorf("Expected status code %d, but got %d", expectedStatusCode, w.Code)
	}
}

func TestCreateItem500PostgresError(t *testing.T) {
	// setup mock dependencies and DB query expectations
	deps, mockDBPool := getMockDependencies()
	mockCreateRecord := mockRecords[mockRecord1]
	mockDBPool.ExpectQuery("INSERT INTO item (.+) VALUES (.+) RETURNING id").
		WithArgs(mockCreateRecord.Name, mockCreateRecord.Price).
		WillReturnError(&pgconn.PgError{Code: "12345"})
	// setup router
	r := gin.Default()
	r.POST("/api/items", routes.HandleCreateItem(deps))
	// exec request
	createItemRequest := models.CreateItemRequest{
		Data: models.ItemIn{
			Name:  mockCreateRecord.Name,
			Price: mockCreateRecord.Price,
		},
	}
	createItemRequestJson, _ := json.Marshal(createItemRequest)
	w := performRequest(r, "POST", "/api/items", string(createItemRequestJson))
	// assert response code
	expectedStatusCode := http.StatusInternalServerError
	if w.Code != expectedStatusCode {
		t.Errorf("Expected status code %d, but got %d", expectedStatusCode, w.Code)
	}
}
