import logging

import asyncpg
from fastapi import APIRouter, Depends, HTTPException, Path, Query

from src import models
from src.database import Database, get_database

logger = logging.getLogger(__name__)


router: APIRouter = APIRouter(
    prefix="/api/items", tags=["items"], dependencies=[Depends(get_database)]
)


@router.get(
    "/{item_id}",
    description="Fetch single item by id.",
    responses={
        "404": {"description": "Resource not found"},
    },
    tags=["items"],
)
async def get_item(
    item_id: int = Path(gt=0, example=1), db: Database = Depends(get_database)
) -> models.ItemOutput_GET:
    logger.info("Fetching item by id", extra={"item_id": item_id})

    FETCH_ITEM_BY_ID_COMMAND = """
        SELECT * FROM item
        WHERE id = $1
    """

    async with db.pool.acquire() as con:
        item_record = await con.fetchrow(FETCH_ITEM_BY_ID_COMMAND, item_id)
        if item_record is None:
            raise HTTPException(status_code=404, detail="Item resource not found")
        item = models.Item(**item_record)
        logger.info("Item fetched", extra={"item": dict(item)})
        return models.ItemOutput_GET(item=item)


@router.get("", description="Fetch multiple items by ids.", tags=["items"])
async def get_items(
    item_ids: list[int] = Query(gt=0, example=[1, 2]), db: Database = Depends(get_database)
) -> models.ItemsOutput_GET:
    logger.info("Fetching items by ids", extra={"item_ids": item_ids})

    FETCH_ITEMS_BY_IDS_COMMAND = """
        SELECT * FROM item
        WHERE id = ANY($1::int[])
    """

    async with db.pool.acquire() as con:
        fetched_records = await con.fetch(FETCH_ITEMS_BY_IDS_COMMAND, item_ids)
        items = [models.Item(**r) for r in fetched_records]
        logger.info("Items fetched", extra={"items": [dict(item) for item in items]})
        return models.ItemsOutput_GET(items=items)


@router.post(
    "",
    description="Save new items.",
    responses={
        "409": {"description": "Resource already exists"},
    },
    tags=["items"],
)
async def create_items(
    input: models.ItemsInput_POST, db: Database = Depends(get_database)
) -> models.ItemsOutput_POST:
    logger.info("Inserting items", extra={"items": input.items})

    INSERT_ITEMS_COMMAND = """
        INSERT INTO item (name, price)
        VALUES ($1, $2)
        RETURNING id;
    """

    FETCH_ITEMS_BY_IDS_COMMAND = """
        SELECT * FROM item
        WHERE id = ANY($1::int[])
    """

    try:
        async with db.pool.acquire() as con:
            created_ids = []
            for input_item in input.items:
                record_id = await con.fetchval(
                    INSERT_ITEMS_COMMAND, input_item.name, input_item.price
                )
                logger.debug("Item record inserted", extra={"record_id": record_id})
                created_ids.append(record_id)
            created_records = await con.fetch(FETCH_ITEMS_BY_IDS_COMMAND, created_ids)

        items = [models.Item(**r) for r in created_records]
        logger.info("Created items", extra={"items": [dict(item) for item in items]})
        return models.ItemsOutput_POST(items_created=items)

    except asyncpg.exceptions.UniqueViolationError as e:
        logger.debug(
            "Item record could not be created because it violated a unique constraint",
            extra={"error": e},
        )
        raise HTTPException(status_code=409, detail="Resource already exists")
