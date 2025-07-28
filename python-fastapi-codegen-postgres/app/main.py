import logging

from fastapi import Depends, FastAPI, HTTPException, Path, status

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
    item: models.ItemCreateRequest,
    db: Database = Depends(get_database),
) -> models.ItemCreateResponse:
    item_in = item.data
    logger.info("Handling create item request", extra={"item": item_in.model_dump()})
    try:
        item_created = await items_repo.create_item(db, item_in)
    except items_repo.UniqueViolationError as e:
        logger.exception("Error while creating item", extra={"error": e})
        raise HTTPException(
            status_code=status.HTTP_409_CONFLICT,
            detail=e.message,
        )
    except Exception as e:
        logger.exception("Unexpected error while creating item", extra={"error": e})
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail="Error while creating item",
        )
    logger.info("Created item", extra={"item": item_created.model_dump()})
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
        logger.exception(
            "Unexpected error fetching item by id",
            extra={"item_id": item_id, "error": e},
        )
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail="Error fetching item by id",
        )
    if item is None:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Item resource not found",
        )
    logger.info("Item fetched", extra={"item": item.model_dump()})
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
        logger.exception("Error while updating item", extra={"error": e})
        raise HTTPException(
            status_code=status.HTTP_409_CONFLICT,
            detail=e.message,
        )
    except Exception as e:
        logger.exception(
            "Unexpected error updating item",
            extra={"item_id": item_id, "error": e},
        )
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail="Error updating item",
        )
    if item_updated is None:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Item resource not found",
        )
    logger.info("Item updated", extra={"item": item_updated.model_dump()})
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
            raise HTTPException(
                status_code=status.HTTP_404_NOT_FOUND,
                detail="Item resource not found",
            )
        await items_repo.delete_item(db, item_id)
    except Exception as e:
        logger.exception(
            "Unexpected error deleting item by id",
            extra={"item_id": item_id, "error": e},
        )
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail="Error deleting item by id",
        )
    logger.info("Item deleted", extra={"item": item.model_dump()})
    return models.ItemDeleteResponse(
        data=item,
        meta=models.ItemMeta(item_status=models.ItemStatus.deleted),
    )
