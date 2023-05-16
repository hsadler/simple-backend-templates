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
            name VARCHAR(50),
            price REAL
        );
    """
    db = await get_database()
    async with db.pool.acquire() as con:
        exec_status = await con.execute(CREATE_TABLES_QUERY)
        logger.info("Creating DB tables", extra={"exec_status": exec_status})


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


class ItemIn(BaseModel):
    name: str
    price: float


class Item(ItemIn):
    id: int


class ItemsOutput_GET(BaseModel):
    items: list[Item]


class ItemsInput_POST(BaseModel):
    items: list[ItemIn]


class ItemsOutput_POST(BaseModel):
    items_created: list[Item]


@app.get("/items", description="Fetch all stored items.")
async def get_items(db: Database = Depends(get_database)) -> ItemsOutput_GET:
    logger.info("Fetching all items.")

    FETCH_ITEMS_QUERY = """
        SELECT * FROM item;
    """
    async with db.pool.acquire() as con:
        fetched_records = await con.fetch(FETCH_ITEMS_QUERY)
        items = [Item(id=r["id"], name=r["name"], price=r["price"]) for r in fetched_records]
        return ItemsOutput_GET(items=items)


@app.post("/items", description="Insert list of provided items.")
async def create_items(
    input: ItemsInput_POST, db: Database = Depends(get_database)
) -> ItemsOutput_POST:
    logger.info("Inserting items.", extra={"to_create": input.items})

    INSERT_ITEMS_COMMAND = """
        INSERT INTO item (name, price)
        VALUES ($1, $2)
        RETURNING id;
    """

    FETCH_ITEMS_BY_IDS_COMMAND = """
        SELECT * FROM item
        WHERE id = ANY($1::int[])
    """

    async with db.pool.acquire() as con:
        created_ids = []
        for input_item in input.items:
            record_id = await con.fetchval(INSERT_ITEMS_COMMAND, input_item.name, input_item.price)
            created_ids.append(record_id)
        created_records = await con.fetch(FETCH_ITEMS_BY_IDS_COMMAND, created_ids)

        logger.info("Created item records.", extra={"item_records": str(created_records)})

        items = [Item(id=r["id"], name=r["name"], price=r["price"]) for r in created_records]
        return ItemsOutput_POST(items_created=items)


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
