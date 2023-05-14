import logging

from fastapi import FastAPI
from pydantic import BaseModel

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = FastAPI()


# http GET http://localhost:8000/status/
@app.get("/status/")
async def status():
    logger.info("Request to /status/")
    return "ok"


class Item(BaseModel):
    name: str
    price: float


# http POST http://localhost:8000/items/ name=apple price:=1.23
@app.post("/items/")
async def create_item(item: Item):
    logger.info("Request to /items/")
    item_dict = item.dict()
    return item_dict
