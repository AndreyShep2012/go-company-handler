## go-company-handler

REST service to manage companies entities, it supports basic CRUD operations: create, update, delete and get one entity.

[MongoDB](https://www.mongodb.com/docs/) is used as database. Mongo's id is used as uuid for company.
[Fiber](https://docs.gofiber.io/) is used as web framework
[slog](https://go.dev/blog/slog) is used as logger, log level can be controlled from config

### Architecture

Service has layered architecture: handler -> service -> repository

Handler:
 - receives http request
 - parses body
 - implements all validation

Service:
 - receives data from the repository
 - responsible for business logic

Repository:
 - DB layer, stores and gets data directly from the DB

### Authorization

Only authenticated users should have access to create, update and delete companies.

[JWT Token](https://en.wikipedia.org/wiki/JSON_Web_Token) is used for authorization.

To generate test token it is possible to use any suitable web site, for example - https://jwt.io/

To authorize request service checks `Authorization` header, also token should have `Bearer ` prefix, for example:

```
curl -X POST http://localhost:8080/api/v1/companies/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "name": "Example Company",
    "description": "This is an example company.",
    "amount_of_employees": 50,
    "registered": true,
    "type": "Corporations"
  }'
```

### Config

Config file `config.yml` is used as config file, to use it should be in the same working directory as a binary. Also all values can be overwritten by environment variables, in addition any variable has default value

### Deployment

There are several ways to start application:

#### developer mode

It just launch application from `main.go`, to make it work correctly MongoDB should be started somewhere separately and config should have appropriate connection string. To do it just use `make` command

```bash
make dev
```

Second dev option is to build binary and launch binary instead of just using `main.go`. 

```bash
make run
```

#### docker-compose mode

`Docker` and `docker-compose` should be installed on a machine

This mode builds docker image for application and starts `docker-compose` file with preinstalled MongoDB, after start application is ready to use

```bash
make run-docker
```

To stop app:

```bash
make stop-docker
```

### Testing

Service uses unit and integration tests

To run unit tests:

```bash
make test-unit
```

Integration tests can be found in folder `./test/integration`. These tests uses `docker-compose` to start all dependencies like DB. To run:

```bash
make test-itg
```

To run all tests, unit and integration and check the coverage:

```bash
make test-all
```

### API

Each company is defined by the following attributes:
 - `ID` (uuid) required - Mongos's id is used
 - `Name` (15 characters) required - unique
 - `Description` (3000 characters) optional
 - `Amount of Employees` (int) required
 - `Registered` (boolean) required
 - `Type` (Corporations | NonProfit | Cooperative | Sole Proprietorship) required

#### Create

Endpoint: `POST /companies/create`

Example:

```bash
curl -X POST http://localhost:8080/api/v1/companies/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "name": "Example Company",
    "description": "This is an example company.",
    "amount_of_employees": 50,
    "registered": true,
    "type": "Corporations"
  }'
```

#### Update (id can be used from response of create operation)

Endpoint: `PATCH /companies/:id`

Example:

```bash
curl -X PATCH http://localhost:8080/api/v1/companies/67dd199ad119e40001f9e8b9 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "name": "New Name",
    "description": "This is an example company.",
    "amount_of_employees": 50,
    "registered": true,
    "type": "Corporations"
  }'
```

#### Get

Endpoint: `GET /companies/:id`

Example:

```bash
curl http://localhost:8080/api/v1/companies/67dd199ad119e40001f9e8b9
```

#### Delete

Endpoint: `DELETE /companies/:id`

Example:

```bash
curl -X DELETE http://localhost:8080/api/v1/companies/67dd199ad119e40001f9e8b9 \
    -H "Authorization: Bearer YOUR_TOKEN"
```

#### Health

Endpoint: `GET /health`

Example:

```bash
curl http://localhost:8080/health
```

#### Version

Endpoint: `GET /version`

Example:

```bash
curl http://localhost:8080/version
```

### Linter

[golangci-lint](https://golangci-lint.run/) linter is used, should be installed before usage. To use it just run:

```bash
make lint
```
