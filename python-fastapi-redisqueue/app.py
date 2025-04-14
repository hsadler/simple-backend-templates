from uuid import uuid4
import os

from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import redis
from rq import Queue
import uvicorn

# FastAPI app
app = FastAPI()

# Redis connection and Queue
redis_host = os.getenv("REDIS_HOST", "localhost")
redis_port = int(os.getenv("REDIS_PORT", 6379))
redis_conn = redis.Redis(host=redis_host, port=redis_port, db=0)
q = Queue(connection=redis_conn)

# Database simulation
fake_db = {}

# Pydantic model for User
class User(BaseModel):
    name: str
    email: str

# Background task to simulate user creation
def create_user_task(user_id: str, user_data: dict):
    fake_db[user_id] = user_data

@app.post("/users/")
async def create_user(user: User):
    user_id = str(uuid4())  # generate unique user ID
    task = q.enqueue(create_user_task, user_id, user.model_dump())
    return {"message": "User creation initiated", "task_id": task.get_id()}

@app.get("/users/{user_id}")
async def read_user(user_id: str):
    user = fake_db.get(user_id)
    if user is None:
        raise HTTPException(status_code=404, detail="User not found")
    return user

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)
