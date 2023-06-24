import datetime
import json
import logging
import uuid
from typing import Any, Generator

import pytest
from fastapi.testclient import TestClient
from pytest_mock import MockFixture

from src.database import Database, get_database
from src.main import app
from src.models import Item


@pytest.fixture
def client(mocker: MockFixture) -> Generator[TestClient, None, None]:
    def override_get_db() -> Any:
        return mocker.MagicMock(spec=Database)

    app.dependency_overrides[get_database] = override_get_db
    root_logger = logging.getLogger()
    root_logger.setLevel(logging.INFO)
    yield TestClient(app)
    del app.dependency_overrides[get_database]


def get_mock_item(id: int = 1) -> Item:
    return Item(
        id=id,
        uuid=uuid.UUID("00000000-0000-0000-0000-000000000000"),
        created_at=datetime.datetime(2021, 8, 15, 18, 0),
        name="mock item",
        price=1.99,
    )


@pytest.mark.parametrize(
    "item_id, expected_status_code",
    [
        (1, 200),
        (2, 200),
        (0, 422),
        (-1, 422),
    ],
)
@pytest.mark.asyncio
async def test_get_item_found_status_code(
    client: TestClient, mocker: MockFixture, item_id: int, expected_status_code: int
) -> None:
    mock_item = get_mock_item()
    mocker.patch("src.routers.items.items_repo.fetch_item_by_id", return_value=mock_item)
    response = client.get(f"/api/items/{item_id}")
    assert response.status_code == expected_status_code


@pytest.mark.parametrize(
    "item_id, expected_response",
    [
        (
            1,
            {
                "data": json.loads(get_mock_item(1).json()),
                "meta": {},
            },
        ),
        (
            2,
            {
                "data": json.loads(get_mock_item(2).json()),
                "meta": {},
            },
        ),
    ],
)
@pytest.mark.asyncio
async def test_get_item_found_response_shape(
    client: TestClient, mocker: MockFixture, item_id: int, expected_response: dict[str, Any]
) -> None:
    mock_item = get_mock_item(item_id)
    mocker.patch("src.routers.items.items_repo.fetch_item_by_id", return_value=mock_item)
    response = client.get(f"/api/items/{item_id}")
    assert response.json() == expected_response


@pytest.mark.parametrize(
    "item_id, expected_status_code",
    [(3, 404), (4, 404)],
)
@pytest.mark.asyncio
async def test_get_item_not_found_status_code(
    client: TestClient, mocker: MockFixture, item_id: int, expected_status_code: int
) -> None:
    mock_item = None
    mocker.patch("src.routers.items.items_repo.fetch_item_by_id", return_value=mock_item)
    response = client.get(f"/api/items/{item_id}")
    print(response.json())
    assert response.status_code == expected_status_code


@pytest.mark.parametrize(
    "item_id, expected_status_code",
    [("abc", 422), ("1.01", 422)],
)
@pytest.mark.asyncio
async def test_get_item_malformed_id_status_code(
    client: TestClient, mocker: MockFixture, item_id: str, expected_status_code: int
) -> None:
    mock_item = get_mock_item()
    mocker.patch("src.routers.items.items_repo.fetch_item_by_id", return_value=mock_item)
    response = client.get(f"/api/items/{item_id}")
    assert response.status_code == expected_status_code


@pytest.mark.parametrize(
    "item_ids, expected_status_code",
    [
        ([1, 2], 200),
        ([1, 2, 3], 200),
        ([], 422),
        ([0], 422),
        ([-1], 422),
    ],
)
@pytest.mark.asyncio
async def test_get_items_found_status_code(
    client: TestClient, mocker: MockFixture, item_ids: list[int], expected_status_code: int
) -> None:
    mock_items = [get_mock_item() for _ in item_ids]
    mocker.patch("src.routers.items.items_repo.fetch_items_by_ids", return_value=mock_items)
    response = client.get("/api/items", params={"item_ids": item_ids})
    assert response.status_code == expected_status_code


@pytest.mark.parametrize(
    "item_ids, expected_response",
    [
        (
            [1, 2],
            {
                "data": [
                    json.loads(get_mock_item(1).json()),
                    json.loads(get_mock_item(2).json()),
                ],
                "meta": {},
            },
        ),
        (
            [1, 2, 3],
            {
                "data": [
                    json.loads(get_mock_item(1).json()),
                    json.loads(get_mock_item(2).json()),
                    json.loads(get_mock_item(3).json()),
                ],
                "meta": {},
            },
        ),
    ],
)
@pytest.mark.asyncio
async def test_get_items_found_response_shape(
    client: TestClient, mocker: MockFixture, item_ids: list[int], expected_response: dict[str, Any]
) -> None:
    mock_items = [get_mock_item(id) for id in item_ids]
    mocker.patch("src.routers.items.items_repo.fetch_items_by_ids", return_value=mock_items)
    response = client.get("/api/items", params={"item_ids": item_ids})
    assert response.json() == expected_response


@pytest.mark.asyncio
async def test_get_items_not_found(client: TestClient, mocker: MockFixture) -> None:
    mock_items: list[Item] = []
    mocker.patch("src.routers.items.items_repo.fetch_items_by_ids", return_value=mock_items)
    response = client.get("/api/items", params={"item_ids": [1, 2]})
    assert response.status_code == 200
    assert response.json() == {"data": [], "meta": {}}


@pytest.mark.parametrize(
    "item_ids, expected_status_code",
    [
        ([1, "two"], 422),
        ([[1], 2, 3], 422),
        ([1, 2, {"three": 3}], 422),
    ],
)
@pytest.mark.asyncio
async def test_get_items_malformed_input(
    client: TestClient, mocker: MockFixture, item_ids: list[int], expected_status_code: int
) -> None:
    mock_items: list[Item] = []
    mocker.patch("src.routers.items.items_repo.fetch_items_by_ids", return_value=mock_items)
    response = client.get("/api/items", params={"item_ids": item_ids})
    assert response.status_code == expected_status_code
