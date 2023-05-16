import asyncpg

from src import settings


class Database:
    def __init__(self, pool: asyncpg.Pool) -> None:
        self.pool: asyncpg.Pool = pool


async def create_db() -> Database:
    pool = await asyncpg.create_pool(settings.DATABASE_URL, min_size=5, max_size=10)
    return Database(pool)
