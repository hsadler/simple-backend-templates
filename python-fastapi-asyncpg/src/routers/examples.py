import logging
import httpx

from fastapi import APIRouter, HTTPException, Path, Response

logger = logging.getLogger(__name__)


router: APIRouter = APIRouter(prefix="/api/external_call/weather", tags=["examples"])

@router.get(
    "/{city}",
    description="Fetch weather by city name.",
    responses={},
    tags=["items"],
)
async def get_item(city: str = Path(min_length=1, example="stockholm")) -> Response:
    logger.info("Fetching weather by city", extra={"city": city})
    try:
        async with httpx.AsyncClient() as client:
            res = await client.get(f"https://wttr.in/{city}?format=3")
        logger.info("Weather fetched", extra={"weather": res.text})
    except Exception as e:
        logger.exception("Error fetching weather by city", extra={"city": city, "error": e})
        raise HTTPException(status_code=500)
    return Response(content=res.text, media_type="text/plain")