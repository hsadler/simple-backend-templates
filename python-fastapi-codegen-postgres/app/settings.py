from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    debug: bool = False
    is_prod: bool = False
    database_url: str = ""


settings = Settings()
