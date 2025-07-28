import logging

import asyncpg

from app.database import Database
from app.models import Item, ItemIn

logger = logging.getLogger(__name__)


class UniqueViolationError(Exception):
    """Raised when a database operation violates a unique constraint."""

    def __init__(self, message: str):
        self.message = message
        super().__init__(self.message)


async def create_item(db: Database, input_item: ItemIn) -> Item:
    INSERT_ITEM_COMMAND = """
        INSERT INTO item (name, price)
        VALUES ($1, $2)
        RETURNING id;
    """

    FETCH_ITEM_BY_ID_COMMAND = """
        SELECT * FROM item
        WHERE id = $1
    """

    try:
        con: asyncpg.Connection
        async with db.pool.acquire() as con:
            item_id = await con.fetchval(INSERT_ITEM_COMMAND, input_item.name, input_item.price)
            item_created_record = await con.fetchrow(FETCH_ITEM_BY_ID_COMMAND, item_id)
        return Item(**item_created_record)

    except asyncpg.exceptions.UniqueViolationError:
        raise UniqueViolationError(message="Item violated a unique constraint")


async def fetch_item(db: Database, item_id: int) -> Item | None:
    FETCH_ITEM_BY_ID_COMMAND = """
        SELECT * FROM item
        WHERE id = $1
    """

    con: asyncpg.Connection
    async with db.pool.acquire() as con:
        item_record = await con.fetchrow(FETCH_ITEM_BY_ID_COMMAND, item_id)

    return Item(**item_record) if item_record else None


async def update_item(db: Database, item_id: int, input_item: ItemIn) -> Item | None:
    UPDATE_ITEM_COMMAND = """
        UPDATE item
        SET name = $1, price = $2
        WHERE id = $3
        RETURNING *;
    """

    try:
        con: asyncpg.Connection
        async with db.pool.acquire() as con:
            item_record = await con.fetchrow(
                UPDATE_ITEM_COMMAND,
                input_item.name,
                input_item.price,
                item_id,
            )
        return Item(**item_record) if item_record else None

    except asyncpg.exceptions.UniqueViolationError:
        raise UniqueViolationError(message="Item violated a unique constraint")


async def delete_item(db: Database, item_id: int) -> None:
    DELETE_ITEM_COMMAND = """
        DELETE FROM item
        WHERE id = $1
    """

    con: asyncpg.Connection
    async with db.pool.acquire() as con:
        await con.execute(DELETE_ITEM_COMMAND, item_id)
