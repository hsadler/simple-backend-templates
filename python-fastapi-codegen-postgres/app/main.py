import uuid
from fastapi import FastAPI

from app import models


app = FastAPI(
    title="Python + FastAPI + CodeGen + Postgres",
    description="A simple FastAPI server with a Postgres database. Models are generated from an OpenAPI schema.",
    version="1.0.0"
)


@app.post("/user")
async def create_user(user: models.UserCreate) -> models.UserResponse:
    """
    Create a new user
    """
    return models.UserResponse(
        id=str(uuid.uuid4()),
        name=user.name,
        age=user.age
    )
