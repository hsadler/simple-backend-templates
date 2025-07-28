import logging
from contextlib import asynccontextmanager
from typing import AsyncGenerator

from fastapi import FastAPI
from prometheus_fastapi_instrumentator import Instrumentator

from app import models
from app.database import get_database
from app.log import setup_logging
from app.routers.items import router as items_router
from app.settings import settings

setup_logging()
logger = logging.getLogger(__name__)


@asynccontextmanager
async def lifespan(_: FastAPI) -> AsyncGenerator[None, None]:
    # Startup logic
    if settings.is_prod:
        db = await get_database()
        await db.run_migrations()
    yield
    # Shutdown logic
    db = await get_database()
    await db.cleanup()


app = FastAPI(
    docs_url="/docs",
    title="Example Python FastAPI Server",
    version="0.1.0",
    lifespan=lifespan,
)


@app.get("/status", description='Returns `"ok"` if the server is up', tags=["status"])
async def status() -> models.StatusOutput:
    logger.info("Request to /status")
    return models.StatusOutput(status="ok")


app.include_router(items_router)


Instrumentator().instrument(app).expose(app)
