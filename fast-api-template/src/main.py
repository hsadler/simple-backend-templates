from fastapi import FastAPI
from pydantic import BaseModel

app = FastAPI()

@app.get("/status/")
async def status():
    return "ok"

class Item(BaseModel):
    name: str
    price: float
    is_offer: bool = None

# http POST http://localhost:8000/items/ name=apple price:=1.23 is_offer:=true
@app.post("/items/")
async def create_item(item: Item):
    item_dict = item.dict()
    if item.is_offer:
        item_dict.update({"message": "Special offer!"})
    return item_dict
