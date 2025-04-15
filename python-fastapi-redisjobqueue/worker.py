import json
import time
from typing import Optional
import os

import redis

redis_host = os.getenv("REDIS_HOST", "localhost")
redis_port = int(os.getenv("REDIS_PORT", 6379))
redis_conn = redis.Redis(host=redis_host, port=redis_port, db=0)

def process_add_numbers(input_data: dict) -> float:
    # Simulate some work
    time.sleep(1)  # Simulate processing time
    x = input_data["x"]
    y = input_data["y"]
    return x + y

def update_job_result(job_id: str, result: Optional[float] = None, error: Optional[str] = None):
    job_key = f"job:{job_id}"
    update_dict = {"status": "complete" if result is not None else "failed"}
    if result is not None:
        update_dict["result"] = str(result)
    if error:
        update_dict["error"] = error
    redis_conn.hset(job_key, mapping=update_dict)

def process_job(job_id: str, job_type: str):
    try:
        # Get job details
        job_key = f"job:{job_id}"
        job_data = redis_conn.hgetall(job_key)
        
        if not job_data:
            print(f"Job {job_id} not found or expired")
            return
            
        input_data = json.loads(job_data[b"input_data"].decode())
        
        # Process based on job type
        if job_type == "add_numbers":
            result = process_add_numbers(input_data)
            update_job_result(job_id, result=result)
        else:
            update_job_result(job_id, error=f"Unknown job type: {job_type}")
            
    except Exception as e:
        update_job_result(job_id, error=str(e))

def run_worker():
    print("Worker started...")
    while True:
        try:
            # Read new messages from the stream, blocking until one arrives
            messages = redis_conn.xread(
                {"jobs_stream": "$"}, 
                block=0,
                count=1
            )
            
            for stream_name, stream_messages in messages:
                for message_id, data in stream_messages:
                    job_id = data[b"job_id"].decode()
                    job_type = data[b"type"].decode()
                    print(f"Processing job: {job_id} of type: {job_type}")
                    process_job(job_id, job_type)
                    
                    # Acknowledge/delete the message
                    redis_conn.xdel("jobs_stream", message_id)
                    
        except Exception as e:
            print(f"Worker error: {e}")
            time.sleep(1)

if __name__ == "__main__":
    run_worker()
