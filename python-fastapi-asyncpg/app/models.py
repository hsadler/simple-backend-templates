from datetime import datetime
from typing import Any
from uuid import UUID

from pydantic import BaseModel, Field

# Resource models


class ItemIn(BaseModel):
    name: str = Field(description="Item name.", examples=["foo"], max_length=50)
    price: float = Field(description="Item price.", examples=[3.14], gt=0)


class Item(BaseModel):
    id: int = Field(description="Item id. Autoincremented.", examples=[1], gt=0)
    uuid: UUID = Field(
        description="Item uuid4 identifier.", examples=["123e4567-e89b-12d3-a456-426614174000"]
    )
    created_at: datetime = Field(
        description="Item time created.", examples=["2021-01-01T00:00:00.000000"]
    )
    name: str = Field(description="Item name.", examples=["foo"], max_length=50)
    price: float = Field(description="Item price.", examples=[3.14], gt=0)


# API i/o models


class StatusOutput(BaseModel):
    status: str = Field(description="Status of the service.")


class ItemInput(BaseModel):
    data: ItemIn = Field(description="Item to be created.")


class ItemOutput(BaseModel):
    data: Item = Field(description="Item.")
    meta: dict[str, Any] = Field(description="Metadata about the item.")


class ItemsOutput(BaseModel):
    data: list[Item] = Field(description="Items.")
    meta: dict[str, Any] = Field(description="Metadata about the items.")
