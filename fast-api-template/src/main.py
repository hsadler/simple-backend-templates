import logging

from fastapi import FastAPI
from pydantic import BaseModel, Field
from starlette_exporter import PrometheusMiddleware, handle_metrics

from src.log import setup_logging

setup_logging()
logger = logging.getLogger(__name__)

app = FastAPI(docs_url="/docs", redoc_url=None)
app.add_middleware(
    PrometheusMiddleware,
    group_paths=True,
    prefix="http",
)
app.add_route("/metrics", handle_metrics)


class StatusOutput(BaseModel):
    status: str = Field(description="Status description")


# http GET http://localhost:8000/status
@app.get("/status")
async def status() -> StatusOutput:
    logger.info("Request to /status")
    return StatusOutput(status="ok")


# EXAMPLE ITEMS API


class Item(BaseModel):
    name: str
    price: float


class ItemsGETOutput(BaseModel):
    items: list[Item]


class ItemsPOSTInput(BaseModel):
    items: list[Item]


class ItemsPOSTOutput(BaseModel):
    items_created: list[Item]


items: list[Item] = []


@app.get("/items")
async def get_items() -> ItemsGETOutput:
    logger.info("GET Request to /items", extra={"return_items": items})
    return ItemsGETOutput(items=items)


@app.post("/items")
async def create_item(input: ItemsPOSTInput) -> ItemsPOSTOutput:
    item: Item
    for item in input.items:
        items.append(item)
    logger.info(
        "POST Request to /items",
        extra={"items_created": [item.dict() for item in input.items]},
    )
    return ItemsPOSTOutput(items_created=input.items)
