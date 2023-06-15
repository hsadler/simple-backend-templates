import logging
from typing import Union

import asyncpg

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


async def create_db() -> Database:
    pool = await asyncpg.create_pool(settings.database_url, min_size=5, max_size=10)
    return Database(pool)


async def init_db() -> None:
    INIT_DB_QUERY = """
        CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
        CREATE TABLE IF NOT EXISTS item (
            id SERIAL PRIMARY KEY,
            uuid UUID DEFAULT uuid_generate_v4(),
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            name VARCHAR(50),
            price NUMERIC(10, 2),
            CONSTRAINT name_unique UNIQUE (name)
        );
    """
    db = await get_database()
    async with db.pool.acquire() as con:
        exec_status = await con.execute(INIT_DB_QUERY)
        logger.info("Initializing DB and creating tables", extra={"exec_status": exec_status})
