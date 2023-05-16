import logging
from typing import Union

import asyncpg
from fastapi import Depends, FastAPI
from pydantic import BaseModel, Field
from starlette_exporter import PrometheusMiddleware, handle_metrics

from src import database
from src.log import setup_logging

setup_logging()
logger = logging.getLogger(__name__)

# INIT FAST API

app = FastAPI(
    docs_url="/docs",
    title="Example Python FastAPI Server",
    version="0.1.0",
)
app.add_middleware(
    PrometheusMiddleware,
    group_paths=True,
    prefix="http",
)
app.add_route("/metrics", handle_metrics)


# INIT DB


db: Union[database.Database, None] = None


async def get_database() -> database.Database:
    global db
    if not db:
        db = await database.create_db()
    return db


# REGISTER APP EVENT LISTENERS


@app.on_event("startup")
async def startup_event() -> None:
    global db
    db = None


@app.on_event("shutdown")
async def shutdown_event() -> None:
    global db
    if db:
        await db.pool.close()
        db = None


# STATUS API


class StatusOutput(BaseModel):
    status: str = Field(description="Status description")


@app.get("/status", description="Provides server status.")
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


@app.get("/items", description="Retrieve all stored items.")
async def get_items() -> ItemsGETOutput:
    logger.info("GET Request to /items", extra={"return_items": items})
    return ItemsGETOutput(items=items)


@app.post("/items", description="Store list of provided items.")
async def create_item(input: ItemsPOSTInput) -> ItemsPOSTOutput:
    item: Item
    for item in input.items:
        items.append(item)
    res = ItemsPOSTOutput(items_created=input.items)
    logger.info(
        "POST Request to /items",
        extra={"items_created": res.items_created},
    )
    return res


# TEST POSTGRES DB STUFF


@app.get("/db", description="test db")
async def get_db(db: database.Database = Depends(get_database)) -> list[str]:
    con: asyncpg.pool.PoolAcquireContext
    async with db.pool.acquire() as con:
        records: list[asyncpg.Record] = await con.fetch(
            """
            SELECT table_name
            FROM information_schema.tables;
            """
        )
        return [str(r) for r in records]
