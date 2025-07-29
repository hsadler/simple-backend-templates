from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    database_url: str = ""
    debug: bool = False
    json_logging: bool = False


settings = Settings()
