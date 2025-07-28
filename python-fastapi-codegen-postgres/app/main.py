import logging

from fastapi import Depends, FastAPI, HTTPException, Query, Response
import asyncpg

from app import models
from app.database import Database, get_database
from app.repos import items as items_repo

logger = logging.getLogger(__name__)


app = FastAPI(
    title="Python + FastAPI + CodeGen + Postgres",
    description="A simple FastAPI server with a Postgres database. Models are generated from an OpenAPI schema.",
    version="1.0.0",
)


@app.get("/ping")
async def ping() -> models.PingResponse:
    return models.PingResponse(message="pong")


@app.post("/items")
async def create_item(
    item: models.CreateItemRequest,
    db: Database = Depends(get_database),
) -> models.CreateItemResponse:
    item_in = item.data
    logger.info("Inserting item", extra={"item": dict(item_in)})
    try:
        item_created = await items_repo.create_item(db, item_in)
    except asyncpg.exceptions.UniqueViolationError:
        raise HTTPException(status_code=409, detail="Resource already exists")
    except Exception as e:
        logger.exception("Error while creating item", extra={"error": e})
        raise HTTPException(status_code=500)
    logger.info("Created item", extra={"item": dict(item_created)})
    return models.CreateItemResponse(data=item_created, meta={"created": True})


@app.get("/items")
async def get_item(
    item_id: int = Query(gt=0, examples=[1]),
    db: Database = Depends(get_database),
) -> models.GetItemResponse:
    logger.info("Fetching item by id", extra={"item_id": item_id})
    try:
        item = await items_repo.fetch_item_by_id(db, item_id)
    except Exception as e:
        logger.exception("Error fetching item by id", extra={"item_id": item_id, "error": e})
        raise HTTPException(status_code=500)
    if item is None:
        raise HTTPException(status_code=404, detail="Item resource not found")
    logger.info("Item fetched", extra={"item": dict(item)})
    return models.GetItemResponse(data=item, meta={})


@app.delete("/items")
async def delete_item(
    item_id: int = Query(gt=0, examples=[1]),
    db: Database = Depends(get_database),
) -> Response:
    logger.info("Deleting item by id", extra={"item_id": item_id})
    try:
        # First check if the item exists
        item = await items_repo.fetch_item_by_id(db, item_id)
        if item is None:
            raise HTTPException(status_code=404, detail="Item resource not found")
        # Delete the item
        await items_repo.delete_item(db, item_id)
    except Exception as e:
        logger.exception("Error deleting item by id", extra={"item_id": item_id, "error": e})
        raise HTTPException(status_code=500)
    logger.info("Item deleted", extra={"item_id": item_id})
    return Response(status_code=204)
