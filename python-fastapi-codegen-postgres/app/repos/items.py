import logging
from typing import Union

import asyncpg

from app.database import Database
from app.models import Item, ItemIn

logger = logging.getLogger(__name__)


async def fetch_item_by_id(db: Database, item_id: int) -> Union[Item, None]:
    FETCH_ITEM_BY_ID_COMMAND = """
        SELECT * FROM item
        WHERE id = $1
    """

    con: asyncpg.Connection
    async with db.pool.acquire() as con:
        item_record = await con.fetchrow(FETCH_ITEM_BY_ID_COMMAND, item_id)

    return Item(**item_record) if item_record is not None else None


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
            logger.info("Item record inserted", extra={"item_id": item_id})
            item_created_record = await con.fetchrow(FETCH_ITEM_BY_ID_COMMAND, item_id)

        return Item(**item_created_record)

    except asyncpg.exceptions.UniqueViolationError as e:
        logger.info(
            "Item record could not be created because it violated a unique constraint",
            extra={"error": e},
        )
        raise e


async def delete_item(db: Database, item_id: int) -> None:
    DELETE_ITEM_COMMAND = """
        DELETE FROM item
        WHERE id = $1
    """

    con: asyncpg.Connection
    async with db.pool.acquire() as con:
        await con.execute(DELETE_ITEM_COMMAND, item_id)
