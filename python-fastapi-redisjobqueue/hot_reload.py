import importlib
from typing import Any, Callable
from functools import wraps


def hot_reload(func: Callable) -> Callable:
    """
    Decorator that reloads the module containing the function before each call.
    This allows for hot-reloading of function implementations without restarting the process.
    """

    @wraps(func)
    def wrapper(*args: Any, **kwargs: Any) -> Any:
        module_name = func.__module__
        module = importlib.import_module(module_name)
        importlib.reload(module)
        reloaded_func = getattr(module, func.__name__)
        return reloaded_func(*args, **kwargs)

    return wrapper
