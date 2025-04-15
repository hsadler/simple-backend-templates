from uuid import uuid4
import json
from enum import Enum
from typing import Optional
import os

from fastapi import FastAPI, HTTPException, Query
from pydantic import BaseModel
import redis
import uvicorn

# FastAPI app
app = FastAPI()

# Redis connection
redis_host = os.getenv("REDIS_HOST", "localhost")
redis_port = int(os.getenv("REDIS_PORT", 6379))
redis_conn = redis.Redis(host=redis_host, port=redis_port, db=0)

# Job status enum
class JobStatus(str, Enum):
    PENDING = "pending"
    COMPLETE = "complete"
    FAILED = "failed"

# Models
class Job(BaseModel):
    id: str
    status: JobStatus
    type: str
    input_data: dict
    result: Optional[float] = None
    error: Optional[str] = None

# Job management
def create_job(job_type: str, input_data: dict, ttl_seconds: int = 3600) -> Job:
    job_id = str(uuid4())
    job = Job(
        id=job_id,
        status=JobStatus.PENDING,
        type=job_type,
        input_data=input_data
    )
    
    # Store job details in a Redis hash with TTL
    job_key = f"job:{job_id}"
    redis_conn.hset(job_key, mapping={
        "status": job.status,
        "type": job.type,
        "input_data": json.dumps(job.input_data)
    })
    redis_conn.expire(job_key, ttl_seconds)
    
    # Add job to the processing queue
    redis_conn.xadd("jobs_stream", {
        "job_id": job_id,
        "type": job_type
    })
    
    return job

def get_job(job_id: str) -> Optional[Job]:
    job_key = f"job:{job_id}"
    job_data = redis_conn.hgetall(job_key)
    
    if not job_data:
        return None
        
    job_dict = {
        "id": job_id,
        "status": job_data[b"status"].decode(),
        "type": job_data[b"type"].decode(),
        "input_data": json.loads(job_data[b"input_data"].decode()),
    }
    
    if b"result" in job_data:
        job_dict["result"] = float(job_data[b"result"].decode())
    if b"error" in job_data:
        job_dict["error"] = job_data[b"error"].decode()
        
    return Job(**job_dict)

# API endpoints
@app.post("/add-numbers")
async def add_numbers(
    x: float = Query(..., description="First number to add"), 
    y: float = Query(..., description="Second number to add")
):
    job = create_job(
        job_type="add_numbers",
        input_data={"x": x, "y": y},
        ttl_seconds=60 * 60 # 1 hour TTL
    )
    return {"message": "Addition job created", "job_id": job.id}

@app.get("/add-numbers/{job_id}")
async def get_addition_result(job_id: str):
    job = get_job(job_id)
    if job is None:
        raise HTTPException(status_code=404, detail="Job not found")
    
    if job.status == JobStatus.PENDING:
        return {"status": "pending", "message": "Job is still being processed"}
    
    if job.status == JobStatus.FAILED:
        return {"status": "failed", "error": job.error}
    
    return {
        "status": "complete",
        "result": job.result,
        "input": job.input_data
    }

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)
