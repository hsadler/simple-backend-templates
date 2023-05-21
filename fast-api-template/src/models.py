from datetime import datetime
from uuid import UUID

from pydantic import BaseModel, Field


class StatusOutput(BaseModel):
    status: str = Field(description="Status description.")


class ItemIn(BaseModel):
    name: str = Field(max_length=50, description="Item name.", example="foo")
    price: float = Field(gt=0, description="Item price.", example="3.14")


class Item(BaseModel):
    id: int = Field(gt=0, description="Item id. Autoincremented.", example=1)
    uuid: UUID = Field(description="Item uuid4 identifier.")
    created_at: datetime = Field(description="Item time created.")
    name: str = Field(max_length=50, description="Item name.", example="foo")
    price: float = Field(gt=0, description="Item price.", example="3.14")


class ItemInput(BaseModel):
    item: ItemIn


class ItemOutput(BaseModel):
    item: Item


class ItemsOutput(BaseModel):
    items: list[Item]
