import logging

from fastapi import FastAPI
from transformers import pipeline
from ray import serve

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = FastAPI()

@serve.deployment()
@serve.ingress(app)
class SentimentAnalysisDeployment:
    def __init__(self):
        self.model = pipeline(
            "sentiment-analysis",
            model="distilbert-base-uncased-finetuned-sst-2-english",
            tokenizer="distilbert-base-uncased-finetuned-sst-2-english",
        )

    @app.get("/analyze/{text}")
    async def analyze(self, text: str):
        logger.info(f"INPUT: {text}")
        return {"text": f"Hello {text}"}
        return self.model(text)

serve.run(SentimentAnalysisDeployment.bind(), name="sentiment-analysis")

# resp = requests.post("http://localhost:8000/analyze/I love this product!")
# logger.info(resp.json())
