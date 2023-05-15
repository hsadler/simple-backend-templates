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


class Item(BaseModel):
    name: str
    price: float


# http POST http://localhost:8000/items name=apple price=1.23
@app.post("/items")
async def create_item(item: Item) -> Item:
    logger.info("Request to /items", extra={"returned_item": item.dict()})
    return item
