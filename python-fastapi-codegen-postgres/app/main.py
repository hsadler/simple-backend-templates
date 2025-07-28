import uuid
from datetime import datetime

from fastapi import FastAPI, Query, Response

from app import models

app = FastAPI(
    title="Python + FastAPI + CodeGen + Postgres",
    description="A simple FastAPI server with a Postgres database. Models are generated from an OpenAPI schema.",
    version="1.0.0",
)


@app.get("/ping")
async def ping() -> models.PingResponse:
    return models.PingResponse(message="pong")


@app.post("/items")
async def create_item(item: models.CreateItemRequest) -> models.CreateItemResponse:
    # STUB
    return models.CreateItemResponse(
        data=models.Item(
            id=1,
            uuid=uuid.uuid4(),
            created_at=datetime.now(),
            name=item.data.name,
            price=item.data.price,
        ),
        meta=models.CreateItemResponseMeta(created=True),
    )


@app.get("/items")
async def get_item(item_id: int = Query(gt=0, examples=[1])) -> models.GetItemResponse:
    # STUB
    return models.GetItemResponse(
        data=models.Item(
            id=1, uuid=uuid.uuid4(), created_at=datetime.now(), name="Item 1", price=100
        ),
        meta={},
    )


@app.delete("/items")
async def delete_item(item_id: int = Query(gt=0, examples=[1])) -> Response:
    # STUB
    return Response(status_code=200)
