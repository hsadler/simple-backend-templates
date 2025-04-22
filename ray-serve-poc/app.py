import logging

from fastapi import FastAPI
from pydantic import BaseModel
from transformers import pipeline
from ray import serve
from ray.serve.handle import DeploymentHandle, DeploymentResponse
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = FastAPI()

class InputText(BaseModel):
    text: str

@serve.deployment()
class SentimentAnalysisDeployment:
    def __init__(self):
        self.model = pipeline("sentiment-analysis", model="distilbert-base-uncased-finetuned-sst-2-english")

    async def __call__(self, text: str):
        return self.model(text)

app = SentimentAnalysisDeployment.bind()
handle: DeploymentHandle = serve.run(app, route_prefix="/analyze")

# Test the service
text = "I love this product!"
logger.info(f"Starting sentiment analysis for text: {text}")
response: DeploymentResponse = handle.remote(text)
logger.info(f"Sentiment analysis result: {response.result()}")

# TODO: serve the model using the Ray Serve API
