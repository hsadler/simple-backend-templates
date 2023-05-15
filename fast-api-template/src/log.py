import logging

from pythonjsonlogger import jsonlogger

from src import settings


def setup_logging() -> None:
    root_logger = logging.getLogger()
    handler = logging.StreamHandler()
    handler.setFormatter(
        jsonlogger.JsonFormatter("%(asctime)s %(levelname)s %(message)s")
    )
    root_logger.addHandler(handler)
    level = logging.DEBUG if settings.DEBUG == "true" else logging.INFO
    root_logger.setLevel(level)
    root_logger.info("Set up logging", extra={"level": level})
