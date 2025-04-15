import logging
import sys
from typing import Optional

def setup_logger(name: str, level: Optional[int] = None) -> logging.Logger:
    """
    Set up a logger with consistent formatting and configuration.
    
    Args:
        name: The name of the logger (typically __name__)
        level: Optional logging level. Defaults to INFO if not specified.
    
    Returns:
        A configured logger instance
    """
    logger = logging.getLogger(name)
    
    if level is None:
        level = logging.INFO
    
    logger.setLevel(level)
    
    # Create console handler with formatting
    console_handler = logging.StreamHandler(sys.stdout)
    console_handler.setLevel(level)
    
    # Create formatter
    formatter = logging.Formatter(
        '%(asctime)s - %(name)s - %(levelname)s - %(message)s',
        datefmt='%Y-%m-%d %H:%M:%S'
    )
    
    # Add formatter to handler
    console_handler.setFormatter(formatter)
    
    # Add handler to logger if it doesn't already have handlers
    if not logger.handlers:
        logger.addHandler(console_handler)
    
    return logger
