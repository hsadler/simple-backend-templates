import logging
from datetime import datetime
from typing import Generator
from unittest.mock import AsyncMock, MagicMock, patch
from uuid import UUID

import pytest
from fastapi.testclient import TestClient

from src import models
from src.database import get_database
from src.main import app
from src.repos import items as items_repo

"""
TESTS WRITTEN BY AN AI ASSISTANT
"""

# Disable logging during tests to reduce noise
logging.disable(logging.CRITICAL)


@pytest.fixture
def mock_database() -> MagicMock:
    """Create a mock database object."""
    mock_db = MagicMock()
    mock_db.pool = AsyncMock()
    return mock_db


@pytest.fixture
def client(mock_database: MagicMock) -> Generator[TestClient, None, None]:
    """Create a test client for the FastAPI app with mocked database dependency."""
    # Override the database dependency to prevent real database connections
    app.dependency_overrides[get_database] = lambda: mock_database

    client = TestClient(app)
    yield client

    # Clean up dependency overrides after each test
    app.dependency_overrides.clear()


@pytest.fixture
def sample_item() -> models.Item:
    """Create a sample Item object for testing."""
    return models.Item(
        id=1,
        uuid=UUID("550e8400-e29b-41d4-a716-446655440000"),
        created_at=datetime(2021, 1, 1, 0, 0, 0),
        name="test_item",
        price=10.0,
    )


@pytest.fixture
def sample_item_in() -> models.ItemIn:
    """Create a sample ItemIn object for testing."""
    return models.ItemIn(name="test_item", price=10.0)


class TestPingEndpoint:
    """Tests for the /ping endpoint."""

    def test_ping_success(self, client: TestClient) -> None:
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
    ) -> None:
        """Test successful item creation."""

        # Mock the create_item function
        with patch("src.repos.items.create_item", return_value=sample_item) as mock_create:
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
    ) -> None:
        """Test item creation with unique constraint violation."""

        # Mock the create_item function to raise a UniqueViolationError
        with patch(
            "src.repos.items.create_item",
            side_effect=items_repo.UniqueViolationError("Item violated a unique constraint"),
        ) as mock_create:
            response = client.post("/items", json={"data": sample_item_in.model_dump()})

            assert response.status_code == 409
            response_data = response.json()
            assert "detail" in response_data
            assert "Item violated a unique constraint" in response_data["detail"]

            # Verify the mock was called correctly
            mock_create.assert_called_once()

    @pytest.mark.asyncio
    async def test_create_item_database_error(
        self,
        client: TestClient,
        mock_database: MagicMock,
        sample_item_in: models.ItemIn,
    ) -> None:
        """Test item creation with database error."""

        # Mock the create_item function to raise a generic Exception
        with patch(
            "src.repos.items.create_item",
            side_effect=Exception("database connection failed"),
        ) as mock_create:
            response = client.post("/items", json={"data": sample_item_in.model_dump()})

            assert response.status_code == 500
            response_data = response.json()
            assert "detail" in response_data
            assert "Error while creating item" in response_data["detail"]

            # Verify the mock was called correctly
            mock_create.assert_called_once()

    def test_create_item_invalid_request_body(self, client: TestClient) -> None:
        """Test item creation with invalid request body."""

        response = client.post("/items", json={"invalid": "data"})

        assert response.status_code == 422

    def test_create_item_missing_data_field(self, client: TestClient) -> None:
        """Test item creation with missing data field."""

        response = client.post("/items", json={"name": "test_item", "price": 10.0})

        assert response.status_code == 422

    def test_create_item_invalid_price(self, client: TestClient) -> None:
        """Test item creation with invalid price."""

        response = client.post("/items", json={"data": {"name": "test_item", "price": -5.0}})

        assert response.status_code == 422


class TestGetItemEndpoint:
    """Tests for the GET /items/{item_id} endpoint."""

    @pytest.mark.asyncio
    async def test_get_item_success(
        self,
        client: TestClient,
        mock_database: MagicMock,
        sample_item: models.Item,
    ) -> None:
        """Test successful item retrieval."""

        # Mock the fetch_item function
        with patch("src.repos.items.fetch_item", return_value=sample_item) as mock_fetch:
            response = client.get("/items/1")

            assert response.status_code == 200
            response_data = response.json()
            assert response_data["data"]["id"] == 1
            assert response_data["data"]["name"] == "test_item"
            assert response_data["data"]["price"] == 10.0
            assert response_data["meta"]["item_status"] == "fetched"

            # Verify the mock was called correctly
            mock_fetch.assert_called_once_with(mock_database, 1)

    @pytest.mark.asyncio
    async def test_get_item_not_found(
        self,
        client: TestClient,
        mock_database: MagicMock,
    ) -> None:
        """Test item retrieval when item doesn't exist."""

        # Mock the fetch_item function to return None
        with patch(
            "src.repos.items.fetch_item",
            return_value=None,
        ) as mock_fetch:
            response = client.get("/items/999")

            assert response.status_code == 404
            response_data = response.json()
            assert "detail" in response_data
            assert "Item resource not found" in response_data["detail"]

            # Verify the mock was called correctly
            mock_fetch.assert_called_once_with(mock_database, 999)

    @pytest.mark.asyncio
    async def test_get_item_database_error(
        self,
        client: TestClient,
        mock_database: MagicMock,
    ) -> None:
        """Test item retrieval with database error."""

        # Mock the fetch_item function to raise a generic Exception
        with patch(
            "src.repos.items.fetch_item",
            side_effect=Exception("database connection failed"),
        ) as mock_fetch:
            response = client.get("/items/1")

            assert response.status_code == 500
            response_data = response.json()
            assert "detail" in response_data
            assert "Error fetching item by id" in response_data["detail"]

            # Verify the mock was called correctly
            mock_fetch.assert_called_once_with(mock_database, 1)

    def test_get_item_invalid_id(self, client: TestClient) -> None:
        """Test item retrieval with invalid item ID."""

        response = client.get("/items/0")

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
    ) -> None:
        """Test successful item update."""

        # Mock the update_item function
        with patch("src.repos.items.update_item", return_value=sample_item) as mock_update:
            response = client.patch("/items/1", json={"data": sample_item_in.model_dump()})

            assert response.status_code == 200
            response_data = response.json()
            assert response_data["data"]["id"] == 1
            assert response_data["data"]["name"] == "test_item"
            assert response_data["data"]["price"] == 10.0
            assert response_data["meta"]["item_status"] == "updated"

            # Verify the mock was called correctly
            mock_update.assert_called_once_with(mock_database, 1, sample_item_in)

    @pytest.mark.asyncio
    async def test_update_item_not_found(
        self,
        client: TestClient,
        mock_database: MagicMock,
        sample_item_in: models.ItemIn,
    ) -> None:
        """Test item update when item doesn't exist."""

        # Mock the update_item function to return None
        with patch(
            "src.repos.items.update_item",
            return_value=None,
        ) as mock_update:
            response = client.patch("/items/999", json={"data": sample_item_in.model_dump()})

            assert response.status_code == 404
            response_data = response.json()
            assert "detail" in response_data
            assert "Item resource not found" in response_data["detail"]

            # Verify the mock was called correctly
            mock_update.assert_called_once_with(mock_database, 999, sample_item_in)

    @pytest.mark.asyncio
    async def test_update_item_unique_violation(
        self,
        client: TestClient,
        mock_database: MagicMock,
        sample_item_in: models.ItemIn,
    ) -> None:
        """Test item update with unique constraint violation."""

        # Mock the update_item function to raise a UniqueViolationError
        with patch(
            "src.repos.items.update_item",
            side_effect=items_repo.UniqueViolationError("Item violated a unique constraint"),
        ) as mock_update:
            response = client.patch("/items/1", json={"data": sample_item_in.model_dump()})

            assert response.status_code == 409
            response_data = response.json()
            assert "detail" in response_data
            assert "Item violated a unique constraint" in response_data["detail"]

            # Verify the mock was called correctly
            mock_update.assert_called_once_with(mock_database, 1, sample_item_in)

    @pytest.mark.asyncio
    async def test_update_item_database_error(
        self,
        client: TestClient,
        mock_database: MagicMock,
        sample_item_in: models.ItemIn,
    ) -> None:
        """Test item update with database error."""

        # Mock the update_item function to raise a generic Exception
        with patch(
            "src.repos.items.update_item",
            side_effect=Exception("database connection failed"),
        ) as mock_update:
            response = client.patch("/items/1", json={"data": sample_item_in.model_dump()})

            assert response.status_code == 500
            response_data = response.json()
            assert "detail" in response_data
            assert "Error updating item" in response_data["detail"]

            # Verify the mock was called correctly
            mock_update.assert_called_once_with(mock_database, 1, sample_item_in)

    def test_update_item_invalid_id(self, client: TestClient) -> None:
        """Test item update with invalid item ID."""

        response = client.patch("/items/0", json={"data": {"name": "test", "price": 10.0}})

        assert response.status_code == 422

    def test_update_item_invalid_request_body(self, client: TestClient) -> None:
        """Test item update with invalid request body."""

        response = client.patch("/items/1", json={"invalid": "data"})

        assert response.status_code == 422

    def test_update_item_invalid_price(self, client: TestClient) -> None:
        """Test item update with invalid price."""

        response = client.patch("/items/1", json={"data": {"name": "test_item", "price": -5.0}})

        assert response.status_code == 422


class TestDeleteItemEndpoint:
    """Tests for the DELETE /items/{item_id} endpoint."""

    @pytest.mark.asyncio
    async def test_delete_item_success(
        self,
        client: TestClient,
        mock_database: MagicMock,
        sample_item: models.Item,
    ) -> None:
        """Test successful item deletion."""

        # Mock both fetch_item and delete_item functions
        with (
            patch("src.repos.items.fetch_item", return_value=sample_item) as mock_fetch,
            patch("src.repos.items.delete_item") as mock_delete,
        ):
            response = client.delete("/items/1")

            assert response.status_code == 200
            response_data = response.json()
            assert response_data["data"]["id"] == 1
            assert response_data["data"]["name"] == "test_item"
            assert response_data["data"]["price"] == 10.0
            assert response_data["meta"]["item_status"] == "deleted"

            # Verify the mocks were called correctly
            mock_fetch.assert_called_once_with(mock_database, 1)
            mock_delete.assert_called_once_with(mock_database, 1)

    @pytest.mark.asyncio
    async def test_delete_item_not_found(
        self,
        client: TestClient,
        mock_database: MagicMock,
    ) -> None:
        """Test item deletion when item doesn't exist."""

        # Mock the fetch_item function to return None, and also mock delete_item to avoid errors
        with (
            patch(
                "src.repos.items.fetch_item",
                return_value=None,
            ) as mock_fetch,
            patch("src.repos.items.delete_item") as mock_delete,
        ):
            response = client.delete("/items/999")

            assert response.status_code == 404
            response_data = response.json()
            assert "detail" in response_data
            assert "Item resource not found" in response_data["detail"]

            # Verify the fetch mock was called but delete was not
            mock_fetch.assert_called_once_with(mock_database, 999)
            mock_delete.assert_not_called()

    @pytest.mark.asyncio
    async def test_delete_item_database_error(
        self,
        client: TestClient,
        mock_database: MagicMock,
    ) -> None:
        """Test item deletion with database error."""

        # Mock the fetch_item function to raise a generic Exception
        with patch(
            "src.repos.items.fetch_item",
            side_effect=Exception("database connection failed"),
        ) as mock_fetch:
            response = client.delete("/items/1")

            assert response.status_code == 500
            response_data = response.json()
            assert "detail" in response_data
            assert "Error deleting item by id" in response_data["detail"]

            # Verify the mock was called correctly
            mock_fetch.assert_called_once_with(mock_database, 1)

    def test_delete_item_invalid_id(self, client: TestClient) -> None:
        """Test item deletion with invalid item ID."""

        response = client.delete("/items/0")

        assert response.status_code == 422
