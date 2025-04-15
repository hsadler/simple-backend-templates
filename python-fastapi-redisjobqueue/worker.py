import json
import time
from typing import Optional
import os

import redis
from logger_config import setup_logger
from hot_reload import hot_reload

logger = setup_logger(__name__)

redis_host = os.getenv("REDIS_HOST", "localhost")
redis_port = int(os.getenv("REDIS_PORT", 6379))
redis_conn = redis.Redis(host=redis_host, port=redis_port, db=0)


def process_add_numbers(input_data: dict) -> float:
    # Simulate some work
    time.sleep(1)  # Simulate processing time
    x = input_data["x"]
    y = input_data["y"]
    result = x + y
    logger.info(f"Processed add_numbers: {x} + {y} = {result}")
    return result


def update_job_result(
    job_id: str, result: Optional[float] = None, error: Optional[str] = None
):
    job_key = f"job:{job_id}"
    update_dict = {"status": "complete" if result is not None else "failed"}
    if result is not None:
        update_dict["result"] = str(result)
    if error:
        update_dict["error"] = error
    redis_conn.hset(job_key, mapping=update_dict)
    logger.info(f"Updated job {job_id} with status: {update_dict['status']}")


@hot_reload
def process_job(job_id: str, job_type: str):
    try:
        # Get job details
        job_key = f"job:{job_id}"
        job_data = redis_conn.hgetall(job_key)

        if not job_data:
            logger.error(f"Job {job_id} not found or expired")
            return

        input_data = json.loads(job_data[b"input_data"].decode())
        logger.debug(f"Processing job {job_id} with input data: {input_data}")

        # Process based on job type
        if job_type == "add_numbers":
            result = process_add_numbers(input_data)
            update_job_result(job_id, result=result)
        else:
            error_msg = f"Unknown job type: {job_type}"
            logger.error(f"Job {job_id} failed: {error_msg}")
            update_job_result(job_id, error=error_msg)

    except Exception as e:
        error_msg = str(e)
        logger.error(f"Error processing job {job_id}: {error_msg}", exc_info=True)
        update_job_result(job_id, error=error_msg)


def run_worker():
    logger.info("Worker started...")
    while True:
        try:
            # Read new messages from the stream, blocking until one arrives
            messages = redis_conn.xread({"jobs_stream": "$"}, block=0, count=1)

            for stream_name, stream_messages in messages:
                for message_id, data in stream_messages:
                    job_id = data[b"job_id"].decode()
                    job_type = data[b"type"].decode()
                    logger.info(f"Received job: {job_id} of type: {job_type}")
                    process_job(job_id, job_type)

                    # Acknowledge/delete the message
                    redis_conn.xdel("jobs_stream", message_id)
                    logger.debug(f"Acknowledged message {message_id} for job {job_id}")

        except Exception as e:
            logger.error(f"Worker error: {e}", exc_info=True)
            time.sleep(1)


if __name__ == "__main__":
    run_worker()
