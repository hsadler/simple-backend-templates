import uuid

from locust import HttpUser, tag, task


class User(HttpUser):
    @task
    @tag("status")
    def status(self) -> None:
        self.client.get("/status")

    @task
    @tag("long-running")
    def fibonacci(self) -> None:
        self.client.get("/api/example/long-running/fibonacci/32")

    @task
    @tag("create-item")
    def create_item(self) -> None:
        self.client.post(
            "/api/items", json={"data": {"name": str(uuid.uuid4()), "price": "101.01"}}
        )
