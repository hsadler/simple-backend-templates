import logging
from datetime import datetime
from unittest.mock import AsyncMock, MagicMock, patch
from uuid import UUID

import pytest
from fastapi.testclient import TestClient

from app import models
from app.database import get_database
from app.main import app
from app.repos import items as items_repo

# Disable logging during tests to reduce noise
logging.disable(logging.CRITICAL)


@pytest.fixture
def mock_database():
    """Create a mock database object."""
    mock_db = MagicMock()
    mock_db.pool = AsyncMock()
    return mock_db


@pytest.fixture
def client(mock_database):
    """Create a test client for the FastAPI app with mocked database dependency."""
    # Override the database dependency to prevent real database connections
    app.dependency_overrides[get_database] = lambda: mock_database

    client = TestClient(app)
    yield client

    # Clean up dependency overrides after each test
    app.dependency_overrides.clear()


@pytest.fixture
def sample_item():
    """Create a sample Item object for testing."""
    return models.Item(
        id=1,
        uuid=UUID("550e8400-e29b-41d4-a716-446655440000"),
        created_at=datetime(2021, 1, 1, 0, 0, 0),
        name="test_item",
        price=10.0,
    )


@pytest.fixture
def sample_item_in():
    """Create a sample ItemIn object for testing."""
    return models.ItemIn(name="test_item", price=10.0)


class TestPingEndpoint:
    """Tests for the /ping endpoint."""

    def test_ping_success(self, client: TestClient):
        """Test successful ping response."""
        response = client.get("/ping")

        assert response.status_code == 200
        assert response.json() == {"message": "pong"}


class TestCreateItemEndpoint:
    """Tests for the POST /items endpoint."""

    @pytest.mark.asyncio
    async def test_create_item_success(
        self,
        client: TestClient,
        mock_database: MagicMock,
        sample_item: models.Item,
        sample_item_in: models.ItemIn,
    ):
        """Test successful item creation."""

        # Mock the create_item function
        with patch("app.repos.items.create_item", return_value=sample_item) as mock_create:
            response = client.post("/items", json={"data": sample_item_in.model_dump()})

            assert response.status_code == 200
            response_data = response.json()
            assert response_data["data"]["id"] == 1
            assert response_data["data"]["name"] == "test_item"
            assert response_data["data"]["price"] == 10.0
            assert response_data["meta"]["item_status"] == "created"

            # Verify the mock was called correctly
            mock_create.assert_called_once()

    @pytest.mark.asyncio
    async def test_create_item_unique_violation(
        self,
        client: TestClient,
        mock_database: MagicMock,
        sample_item_in: models.ItemIn,
    ):
        """Test item creation with unique constraint violation."""

        # Mock the create_item function to raise UniqueViolationError
        with patch(
            "app.repos.items.create_item",
            side_effect=items_repo.UniqueViolationError("Item violated a unique constraint"),
        ):
            response = client.post("/items", json={"data": sample_item_in.model_dump()})

            assert response.status_code == 409
            assert "Item violated a unique constraint" in response.json()["detail"]

    @pytest.mark.asyncio
    async def test_create_item_unexpected_error(
        self,
        client: TestClient,
        mock_database: MagicMock,
        sample_item_in: models.ItemIn,
    ):
        """Test item creation with unexpected error."""

        # Mock the create_item function to raise a generic exception
        with patch(
            "app.repos.items.create_item",
            side_effect=Exception("Database connection failed"),
        ):
            response = client.post("/items", json={"data": sample_item_in.model_dump()})

            assert response.status_code == 500
            assert "Error while creating item" in response.json()["detail"]

    def test_create_item_invalid_data(self, client: TestClient):
        """Test item creation with invalid request data."""
        # Test with missing required fields
        response = client.post("/items", json={"data": {"name": "test"}})  # missing price
        assert response.status_code == 422

        # Test with negative price
        response = client.post("/items", json={"data": {"name": "test", "price": -1.0}})
        assert response.status_code == 422


class TestGetItemEndpoint:
    """Tests for the GET /items/{item_id} endpoint."""

    @pytest.mark.asyncio
    async def test_get_item_success(
        self,
        client: TestClient,
        mock_database: MagicMock,
        sample_item: models.Item,
    ):
        """Test successful item retrieval."""

        with patch("app.repos.items.fetch_item", return_value=sample_item) as mock_fetch:
            response = client.get("/items/1")

            assert response.status_code == 200
            response_data = response.json()
            assert response_data["data"]["id"] == 1
            assert response_data["data"]["name"] == "test_item"
            assert response_data["data"]["price"] == 10.0
            assert response_data["meta"]["item_status"] == "fetched"

            mock_fetch.assert_called_once_with(mock_database, 1)

    @pytest.mark.asyncio
    async def test_get_item_not_found(
        self,
        client: TestClient,
        mock_database: MagicMock,
    ):
        """Test item retrieval when item doesn't exist."""

        with patch(
            "app.repos.items.fetch_item",
            return_value=None,
        ):
            response = client.get("/items/999")

            assert response.status_code == 404
            assert "Item resource not found" in response.json()["detail"]

    @pytest.mark.asyncio
    async def test_get_item_unexpected_error(
        self,
        client: TestClient,
        mock_database: MagicMock,
    ):
        """Test item retrieval with unexpected error."""

        with patch(
            "app.repos.items.fetch_item",
            side_effect=Exception("Database connection failed"),
        ):
            response = client.get("/items/1")

            assert response.status_code == 500
            assert "Error fetching item by id" in response.json()["detail"]

    def test_get_item_invalid_id(self, client: TestClient):
        """Test item retrieval with invalid item ID."""
        # Test with non-positive ID
        response = client.get("/items/0")
        assert response.status_code == 422

        response = client.get("/items/-1")
        assert response.status_code == 422


class TestUpdateItemEndpoint:
    """Tests for the PATCH /items/{item_id} endpoint."""

    @pytest.mark.asyncio
    async def test_update_item_success(
        self,
        client: TestClient,
        mock_database: MagicMock,
        sample_item: models.Item,
        sample_item_in: models.ItemIn,
    ):
        """Test successful item update."""

        updated_item = models.Item(
            id=1,
            uuid=sample_item.uuid,
            created_at=sample_item.created_at,
            name="updated_item",
            price=20.0,
        )

        with patch("app.repos.items.update_item", return_value=updated_item) as mock_update:
            response = client.patch(
                "/items/1", json={"data": {"name": "updated_item", "price": 20.0}}
            )

            assert response.status_code == 200
            response_data = response.json()
            assert response_data["data"]["id"] == 1
            assert response_data["data"]["name"] == "updated_item"
            assert response_data["data"]["price"] == 20.0
            assert response_data["meta"]["item_status"] == "updated"

            mock_update.assert_called_once()

    @pytest.mark.asyncio
    async def test_update_item_not_found(
        self,
        client: TestClient,
        mock_database: MagicMock,
        sample_item_in: models.ItemIn,
    ):
        """Test item update when item doesn't exist."""

        with patch("app.repos.items.update_item", return_value=None):
            response = client.patch("/items/999", json={"data": sample_item_in.model_dump()})

            assert response.status_code == 404
            assert "Item resource not found" in response.json()["detail"]

    @pytest.mark.asyncio
    async def test_update_item_unique_violation(
        self,
        client: TestClient,
        mock_database: MagicMock,
        sample_item_in: models.ItemIn,
    ):
        """Test item update with unique constraint violation."""

        with patch(
            "app.repos.items.update_item",
            side_effect=items_repo.UniqueViolationError("Item violated a unique constraint"),
        ):
            response = client.patch("/items/1", json={"data": sample_item_in.model_dump()})

            assert response.status_code == 409
            assert "Item violated a unique constraint" in response.json()["detail"]

    @pytest.mark.asyncio
    async def test_update_item_unexpected_error(
        self,
        client: TestClient,
        mock_database: MagicMock,
        sample_item_in: models.ItemIn,
    ):
        """Test item update with unexpected error."""

        with patch(
            "app.repos.items.update_item",
            side_effect=Exception("Database connection failed"),
        ):
            response = client.patch("/items/1", json={"data": sample_item_in.model_dump()})

            assert response.status_code == 500
            assert "Error updating item" in response.json()["detail"]

    def test_update_item_invalid_data(self, client: TestClient):
        """Test item update with invalid request data."""
        # Test with missing required fields
        response = client.patch("/items/1", json={"data": {"name": "test"}})  # missing price
        assert response.status_code == 422

        # Test with negative price
        response = client.patch("/items/1", json={"data": {"name": "test", "price": -1.0}})
        assert response.status_code == 422

    def test_update_item_invalid_id(self, client: TestClient):
        """Test item update with invalid item ID."""
        response = client.patch("/items/0", json={"data": {"name": "test", "price": 10.0}})
        assert response.status_code == 422


class TestDeleteItemEndpoint:
    """Tests for the DELETE /items/{item_id} endpoint."""

    @pytest.mark.asyncio
    async def test_delete_item_success(
        self,
        client: TestClient,
        mock_database: MagicMock,
        sample_item: models.Item,
    ):
        """Test successful item deletion."""

        with (
            patch(
                "app.repos.items.fetch_item",
                return_value=sample_item,
            ) as mock_fetch,
            patch(
                "app.repos.items.delete_item",
                return_value=None,
            ) as mock_delete,
        ):

            response = client.delete("/items/1")

            assert response.status_code == 200
            response_data = response.json()
            assert response_data["data"]["id"] == 1
            assert response_data["data"]["name"] == "test_item"
            assert response_data["meta"]["item_status"] == "deleted"

            mock_fetch.assert_called_once_with(mock_database, 1)
            mock_delete.assert_called_once_with(mock_database, 1)

    @pytest.mark.asyncio
    async def test_delete_item_not_found(
        self,
        client: TestClient,
        mock_database: MagicMock,
    ):
        """Test item deletion when item doesn't exist."""

        with patch("app.repos.items.fetch_item", return_value=None):
            response = client.delete("/items/999")

            assert response.status_code == 404
            assert "Item resource not found" in response.json()["detail"]

    @pytest.mark.asyncio
    async def test_delete_item_unexpected_error(
        self,
        client: TestClient,
        mock_database: MagicMock,
        sample_item: models.Item,
    ):
        """Test item deletion with unexpected error."""

        with patch(
            "app.repos.items.fetch_item",
            side_effect=Exception("Database connection failed"),
        ):
            response = client.delete("/items/1")

            assert response.status_code == 500
            assert "Error deleting item by id" in response.json()["detail"]

    def test_delete_item_invalid_id(self, client: TestClient):
        """Test item deletion with invalid item ID."""
        response = client.delete("/items/0")
        assert response.status_code == 422

        response = client.delete("/items/-1")
        assert response.status_code == 422
