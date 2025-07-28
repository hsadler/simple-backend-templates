# simple-backend-templates

Simple templates for backends. All implement the same "items" example CRUD API.

- [Golang + ogen + pgx](./golang-ogen-pgx/)
- [Golang + Gin + pgx](./golang-gin-pgx/)
- [Python + FastAPI + CodeGen + Postgres](./python-fastapi-codegen-postgres/)
- [Python + FastAPI + asyncpg](./python-fastapi-asyncpg/)

POCs:

- [Python + FastAPI + Redis Job Queue](./python-fastapi-redisjobqueue/)

Tooling examples:

- [uv example](./uv-example/) - A simple example of a Python project using uv.


## Development

Requirements:
- pre-commit

1. Install pre-commit hooks.
```sh
pre-commit install
```

2. Install dev dependencies for all template projects so that pre-commit can run properly.
