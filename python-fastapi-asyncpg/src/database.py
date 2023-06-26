import logging
from typing import Union

import asyncpg

from alembic import command
from alembic.config import Config
from src.settings import settings

logger = logging.getLogger(__name__)


class Database:
    def __init__(self, pool: asyncpg.Pool) -> None:
        self.pool: asyncpg.Pool = pool

    async def cleanup(self) -> None:
        await self.pool.close()


db: Union[Database, None] = None


async def get_database() -> Database:
    global db
    if not db:
        db = await create_db()
        logger.info("Database with new connection pool created")
    return db


async def run_migrations() -> None:
    logger.info("Running DB migrations")
    alembic_cfg = Config("alembic.ini")
    alembic_cfg.set_main_option("sqlalchemy.url", settings.database_url)
    command.upgrade(alembic_cfg, "head")


async def create_db() -> Database:
    pool = await asyncpg.create_pool(settings.database_url, min_size=5, max_size=10)
    return Database(pool)
