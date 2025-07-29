import logging
import sys

from pythonjsonlogger import json as jsonlogger

from src.settings import settings

BASE_LOG_FORMAT = "%(asctime)s %(levelname)s %(name)s %(message)s"
UVICORN_LOGGERS = ["uvicorn", "uvicorn.access"]


def setup_root_json_logging() -> None:
    root_logger = logging.getLogger()
    log_fmt: logging.Formatter
    if settings.json_logging:
        log_fmt = jsonlogger.JsonFormatter(BASE_LOG_FORMAT)
    else:
        log_fmt = logging.Formatter(BASE_LOG_FORMAT)

    # handler for INFO and DEBUG logs to stdout
    stdout_handler = logging.StreamHandler(sys.stdout)
    stdout_handler.setFormatter(log_fmt)
    # level for stdout logs
    stdout_handler.setLevel(logging.DEBUG)
    # stout should only receive INFO and DEBUG
    stdout_handler.addFilter(lambda record: record.levelno <= logging.INFO)
    root_logger.addHandler(stdout_handler)

    # handler for WARNING, ERROR, and CRITICAL to stderr
    stderr_handler = logging.StreamHandler(sys.stderr)
    stderr_handler.setFormatter(log_fmt)
    # level for stderr logs
    stderr_handler.setLevel(logging.WARNING)
    root_logger.addHandler(stderr_handler)

    # log level for the root logger (will apply to both handlers)
    level = logging.DEBUG if settings.debug else logging.INFO
    root_logger.setLevel(level)
    root_logger.info("Set up logging", extra={"level": level})


def setup_uvicorn_json_logging() -> None:
    # Force the uvicorn loggers to inherit the (json) root-logger handlers
    for logger_name in UVICORN_LOGGERS:
        logger = logging.getLogger(logger_name)
        for handler in logger.handlers:
            logger.removeHandler(handler)
        logger.propagate = True


def setup_logging() -> None:
    setup_root_json_logging()
    setup_uvicorn_json_logging()
