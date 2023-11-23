# Hero Ticket Server

## Table of Contents

- [Prerequisites](#prerequisites)
- [Run](#run)
- [Swagger API Documentation](#swagger-api-documentation)
    - [Generate Swagger API Documentation](#generate-swagger-api-documentation)
    - [Run Swagger API Documentation Server](#run-swagger-api-documentation-server)
    - [Swagger API Documentation URL](#swagger-api-documentation-url)

## Prerequisites

- [Go 1.21](https://golang.org/doc/install)
- [Docker](https://docs.docker.com/get-docker/)
- [Make](https://www.gnu.org/software/make/)

##  Run

```bash
$ make up
```

## Swagger API Documentation

### Generate Swagger API Documentation

```bash
$ make swag_gen
```

### Run Swagger API Documentation Server

```bash
$ make swagger
```

### Swagger API Documentation URL

```
http://localhost:1323/swagger/index.html
```