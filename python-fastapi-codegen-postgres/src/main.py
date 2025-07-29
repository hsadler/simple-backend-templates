import logging

from fastapi import Depends, FastAPI, HTTPException, Path, status

from src import models
from src.database import Database, get_database
from src.log import setup_logging
from src.repos import items as items_repo

setup_logging()
logger = logging.getLogger(__name__)


app = FastAPI(
    title="Python + FastAPI + CodeGen + Postgres",
    description=(
        "A simple FastAPI server with a Postgres database. "
        "Models are generated from an OpenAPI schema."
    ),
    version="1.0.0",
)


@app.get("/ping")
async def ping() -> models.PingResponse:
    logger.info("Handling ping request")
    return models.PingResponse(message="pong")


@app.post("/items")
async def create_item(
    item: models.ItemCreateRequest,
    db: Database = Depends(get_database),
) -> models.ItemCreateResponse:
    item_in = item.data
    logger.info("Handling create item request", extra={"item": item_in.model_dump()})
    try:
        item_created = await items_repo.create_item(db, item_in)
    except items_repo.UniqueViolationError as e:
        logger.warning("Error while creating item", extra={"error": e})
        raise HTTPException(
            status_code=status.HTTP_409_CONFLICT,
            detail=e.message,
        )
    except Exception as e:
        logger.error(
            "Unexpected error while creating item",
            extra={"error": e, "item_in": item_in.model_dump()},
        )
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail="Error while creating item",
        )
    logger.debug(f"Created item: {item_created.model_dump()}")
    return models.ItemCreateResponse(
        data=item_created,
        meta=models.ItemMeta(item_status=models.ItemStatus.created),
    )


@app.get("/items/{item_id}")
async def get_item(
    item_id: int = Path(gt=0, examples=[1]),
    db: Database = Depends(get_database),
) -> models.ItemGetResponse:
    logger.info("Handling get item request", extra={"item_id": item_id})
    try:
        item = await items_repo.fetch_item(db, item_id)
    except Exception as e:
        logger.error(
            "Unexpected error fetching item by id",
            extra={"error": e, "item_id": item_id},
        )
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail="Error fetching item by id",
        )
    if item is None:
        logger.debug(f"Item not found for fetch in db by item_id: {item_id}")
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Item resource not found",
        )
    logger.debug(f"Item fetched: {item.model_dump()}")
    return models.ItemGetResponse(
        data=item,
        meta=models.ItemMeta(item_status=models.ItemStatus.fetched),
    )


@app.patch("/items/{item_id}")
async def update_item(
    item: models.ItemUpdateRequest,
    item_id: int = Path(gt=0, examples=[1]),
    db: Database = Depends(get_database),
) -> models.ItemUpdateResponse:
    item_in = item.data
    logger.info("Handling update item request", extra={"item": item_in.model_dump()})
    try:
        item_updated = await items_repo.update_item(db, item_id, item_in)
    except items_repo.UniqueViolationError as e:
        logger.warning("Error while updating item", extra={"error": e})
        raise HTTPException(
            status_code=status.HTTP_409_CONFLICT,
            detail=e.message,
        )
    except Exception as e:
        logger.error(
            "Unexpected error updating item",
            extra={"error": e, "item_id": item_id, "item_in": item_in.model_dump()},
        )
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail="Error updating item",
        )
    if item_updated is None:
        logger.debug(f"Item not found for update in db by item_id: {item_id}")
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Item resource not found",
        )
    logger.debug(f"Item updated: {item_updated.model_dump()}")
    return models.ItemUpdateResponse(
        data=item_updated,
        meta=models.ItemMeta(item_status=models.ItemStatus.updated),
    )


@app.delete("/items/{item_id}")
async def delete_item(
    item_id: int = Path(gt=0, examples=[1]),
    db: Database = Depends(get_database),
) -> models.ItemDeleteResponse:
    logger.info("Handling delete item request", extra={"item_id": item_id})
    try:
        item = await items_repo.fetch_item(db, item_id)
        if item is None:
            logger.debug("Item not found for delete", extra={"item_id": item_id})
            raise HTTPException(
                status_code=status.HTTP_404_NOT_FOUND,
                detail="Item resource not found",
            )
        await items_repo.delete_item(db, item_id)
    except HTTPException:
        # Re-raise HTTPExceptions (like 404) so they don't get caught by the general handler
        raise
    except Exception as e:
        logger.error(
            "Unexpected error deleting item by id",
            extra={"error": e, "item_id": item_id},
        )
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail="Error deleting item by id",
        )
    logger.debug(f"Item deleted: {item.model_dump()}")
    return models.ItemDeleteResponse(
        data=item,
        meta=models.ItemMeta(item_status=models.ItemStatus.deleted),
    )
