import requests
import time
from typing import Any
import fire

# API configuration
API_BASE_URL = "http://localhost:8000"


class JobClient:
    """Client for interacting with the number addition job queue service."""

    def __init__(self, poll_interval: float = 0.5, max_attempts: int = 60):
        """
        Initialize the job client.

        Args:
            poll_interval: Time between polling attempts in seconds
            max_attempts: Maximum number of polling attempts
        """
        self.poll_interval = poll_interval
        self.max_attempts = max_attempts

    def _submit_addition_job(self, x: float, y: float) -> str:
        """Submit a new addition job and return the job ID."""
        response = requests.post(f"{API_BASE_URL}/add-numbers", params={"x": x, "y": y})
        response.raise_for_status()
        data = response.json()
        return data["job_id"]

    def _get_job_result(self, job_id: str) -> dict[str, Any]:
        """Get the result of a job by its ID."""
        response = requests.get(f"{API_BASE_URL}/add-numbers/{job_id}")
        response.raise_for_status()
        return response.json()

    def _poll_job_result(self, job_id: str) -> dict[str, Any]:
        """
        Poll for job results until completion or failure.

        Args:
            job_id: The ID of the job to poll

        Returns:
            The final job result

        Raises:
            TimeoutError: If the job doesn't complete within max_attempts
            RuntimeError: If the job fails
        """
        attempt = 0
        while attempt < self.max_attempts:
            result = self._get_job_result(job_id)

            if result["status"] == "complete":
                return result
            elif result["status"] == "failed":
                raise RuntimeError(
                    f"Job failed: {result.get('error', 'Unknown error')}"
                )

            print(
                f"Job {job_id} is still pending... (attempt {attempt + 1}/{self.max_attempts})"
            )
            time.sleep(self.poll_interval)
            attempt += 1

        raise TimeoutError(
            f"Job {job_id} did not complete within {self.max_attempts} attempts"
        )

    def add(self, x: float, y: float) -> dict[str, Any]:
        """
        Submit a job to add two numbers and wait for the result.

        Args:
            x: First number to add
            y: Second number to add

        Returns:
            Dictionary containing the job result

        Example:
            To add 5 and 3:
            $ python client.py add --x=5 --y=3
        """
        try:
            print(f"Submitting addition job: {x} + {y}")
            job_id = self._submit_addition_job(x, y)
            print(f"Job submitted with ID: {job_id}")

            print("Polling for results...")
            result = self._poll_job_result(job_id)

            print("\nJob completed successfully!")
            print(f"Input: {result['input']}")
            print(f"Result: {result['result']}")

            return result

        except requests.exceptions.RequestException as e:
            print(f"Error communicating with API: {e}")
            raise
        except (TimeoutError, RuntimeError) as e:
            print(f"Error: {e}")
            raise
        except Exception as e:
            print(f"Unexpected error: {e}")
            raise


if __name__ == "__main__":
    fire.Fire(JobClient())
