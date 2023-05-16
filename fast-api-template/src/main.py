import logging
from typing import Union

import asyncpg
from fastapi import Depends, FastAPI
from pydantic import BaseModel, Field
from starlette_exporter import PrometheusMiddleware, handle_metrics

from src.database import Database, create_db
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


# INIT DB AND TABLES


db: Union[Database, None] = None


async def get_database() -> Database:
    global db
    if not db:
        db = await create_db()
        logger.info("Database with new connection pool created")
    return db


async def create_tables() -> None:
    CREATE_TABLES_QUERY = """
        CREATE TABLE IF NOT EXISTS item (
            id SERIAL PRIMARY KEY,
            item_data JSONB
        );
    """
    db = await get_database()
    async with db.pool.acquire() as con:
        status = await con.execute(CREATE_TABLES_QUERY)
        logger.info("Creating DB tables", extra={"status": status})


@app.on_event("startup")
async def startup() -> None:
    await create_tables()


@app.on_event("shutdown")
async def shutdown() -> None:
    global db
    if db:
        await db.cleanup()
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


class ItemsOutput_GET(BaseModel):
    items: list[Item]


class ItemsInput_POST(BaseModel):
    items: list[Item]


class ItemsOutput_POST(BaseModel):
    items_created: list[Item]


items: list[Item] = []


@app.get("/items", description="Retrieve all stored items.")
async def get_items() -> ItemsOutput_GET:
    logger.info("GET Request to /items", extra={"return_items": items})
    return ItemsOutput_GET(items=items)


@app.post("/items", description="Store list of provided items.")
async def create_item(input: ItemsInput_POST) -> ItemsOutput_POST:
    item: Item
    for item in input.items:
        items.append(item)
    res = ItemsOutput_POST(items_created=input.items)
    logger.info(
        "POST Request to /items",
        extra={"items_created": res.items_created},
    )
    return res


# TEST POSTGRES DB STUFF


@app.get("/db-tables", description="Show created DB tables")
async def get_db_tables(db: Database = Depends(get_database)) -> list[str]:
    async with db.pool.acquire() as con:
        records: list[asyncpg.Record] = await con.fetch(
            """
            SELECT table_name
            FROM information_schema.tables
            WHERE table_schema = 'public'
            AND table_type = 'BASE TABLE';
        """
        )
        return [str(r) for r in records]
