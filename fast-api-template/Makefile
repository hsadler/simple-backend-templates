up:
	poetry run uvicorn src.main:app --reload

black:
	poetry run black .

isort:
	poetry run isort .

flake8:
	poetry run flake8 .

mypy:
	poetry run mypy .

lint: black flake8 isort mypy
