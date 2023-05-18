import logging

import asyncpg
from fastapi import APIRouter, Depends, HTTPException, Path, Query

from src import models
from src.database import Database, get_database
from src.repos import items as items_repo

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
    item = await items_repo.fetch_item_by_id(db, item_id)
    if item is None:
        raise HTTPException(status_code=404, detail="Item resource not found")
    logger.info("Item fetched", extra={"item": dict(item)})
    return models.ItemOutput_GET(item=item)


@router.get("", description="Fetch multiple items by ids.", tags=["items"])
async def get_items(
    item_ids: list[int] = Query(gt=0, example=[1, 2]), db: Database = Depends(get_database)
) -> models.ItemsOutput_GET:
    logger.info("Fetching items by ids", extra={"item_ids": item_ids})
    items = await items_repo.fetch_items_by_ids(db, item_ids)
    logger.info("Items fetched", extra={"items": [dict(item) for item in items]})
    return models.ItemsOutput_GET(items=items)


@router.post(
    "",
    description="Save new item.",
    responses={
        "409": {"description": "Resource already exists"},
    },
    tags=["items"],
)
async def create_item(
    input: models.ItemInput_POST, db: Database = Depends(get_database)
) -> models.ItemOutput_POST:
    logger.info("Inserting item", extra={"item": dict(input.item)})

    try:
        item_created = await items_repo.create_item(db, input.item)
    except asyncpg.exceptions.UniqueViolationError:
        raise HTTPException(status_code=409, detail="Resource already exists")

    logger.info("Created item", extra={"item": dict(item_created)})
    return models.ItemOutput_POST(item_created=item_created)
