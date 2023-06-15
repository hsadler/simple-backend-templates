from pydantic import BaseSettings


class Settings(BaseSettings):
    debug: bool = False
    database_url: str = ""


settings = Settings()
