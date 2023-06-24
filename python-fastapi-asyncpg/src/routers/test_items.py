import datetime
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
async def test_get_item_status_code(
    client: TestClient, mocker: MockFixture, item_id: int, expected_status_code: int
) -> None:
    mock_item = Item(
        id=1,
        uuid=uuid.UUID("00000000-0000-0000-0000-000000000000"),
        created_at=datetime.datetime(2021, 8, 15, 18, 0),
        name="mock item",
        price=1.99,
    )
    mocker.patch("src.routers.items.items_repo.fetch_item_by_id", return_value=mock_item)
    response = client.get(f"/api/items/{item_id}")
    print(response.json())
    assert response.status_code == expected_status_code


@pytest.mark.parametrize(
    "item_id, expected_response",
    [
        (
            1,
            {
                "data": {
                    "id": 1,
                    "uuid": "00000000-0000-0000-0000-000000000000",
                    "created_at": "2021-08-15T18:00:00",
                    "name": "mock item",
                    "price": 1.99,
                },
                "meta": {},
            },
        ),
        (
            2,
            {
                "data": {
                    "id": 2,
                    "uuid": "00000000-0000-0000-0000-000000000000",
                    "created_at": "2021-08-15T18:00:00",
                    "name": "mock item",
                    "price": 1.99,
                },
                "meta": {},
            },
        ),
    ],
)
@pytest.mark.asyncio
async def test_get_item_response_format(
    client: TestClient, mocker: MockFixture, item_id: int, expected_response: dict[str, Any]
) -> None:
    mock_item = Item(
        id=item_id,
        uuid=uuid.UUID("00000000-0000-0000-0000-000000000000"),
        created_at=datetime.datetime(2021, 8, 15, 18, 0),
        name="mock item",
        price=1.99,
    )
    mocker.patch("src.routers.items.items_repo.fetch_item_by_id", return_value=mock_item)
    response = client.get(f"/api/items/{item_id}")
    assert response.json() == expected_response


# @pytest.mark.parametrize(
#     "item_ids, expected_status_code",
#     [
#         ([1, 2], 200),
#         ([1, 2, 3], 200),
#         ([3], 200),
#         ([], 200),
#         ([0], 422),
#         ([-1], 422),
#     ],
# )
# def test_get_items(client: TestClient, item_ids: list[int], expected_status_code: int) -> None:
#     response = client.get("/api/items", params={"item_ids": item_ids})
#     assert response.status_code == expected_status_code

# @pytest.mark.parametrize(
#     "item_in, expected_status_code",
#     [
#         ({"name": "Item 1", "description": "Item 1 description"}, 200),
#         ({"name": "Item 2", "description": "Item 2 description"}, 200),
#         ({"name": "Item 1", "description": "Item 1 description"}, 409),
#     ],
# )
# def test_create_item(
#     client: TestClient,
#     item_in: dict[str, str],
#     expected_status_code: int
# ) -> None:
#     response = client.post("/api/items", json=item_in)
#     assert response.status_code == expected_status_code

# # Run the tests:
# # $ pytest -v
# #
