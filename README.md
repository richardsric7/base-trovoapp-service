# TrovoApp API


## Setup Instructions

Install Go on your device, then follow these instructions below:

- `$ git clone https://github.com/trovotech-technologies/trovo-wallet-api.git` the repo.
- Run `$ go test -v ./...` to ensure everything is running smoothly.
- Create a `.env` file. A sample `.env.example` can be found on the root document.
- Run `$ go run main.go` to serve on `0.0.0.0:8080` (for windows "`localhost:8080`")


## Tech Stack Overview

The APIs are written in Go. Notably, we use [Go Gin](https://github.com/gin-gonic/gin) as our web framework.

Logical components are grouped into folders under the /internal folder.

We strive to use  [REST API guideline standards](https://stackoverflow.blog/2020/03/02/best-practices-for-rest-api-design/) as much as we can. 

TrovoApp API has been tested to work with the following databases:
- Postgres

## Error Messages

TrovoApp API errors are generally structured as follows

|   name   |  type  |                                                                        description                                                                        |
| -------- | ------ | --------------------------------------------------------------------------------------------------------------------------------------------------------- |
| error     | string    | Type of error. |
| message    | string | Human readable description of the error. Localized per query language.                                                                                                                    |
| data   | string | Extra developer data related to the error  |


## i18n

The project at its current stage does not yet support i18n

## Email

For development purposes, use online fake smtp servers like mailtrap.io


# Setup Documentation

## Required Dependencies

Before proceeding with the setup, ensure that the following dependencies are installed on your system:

1. **CockroachDB Standalone Single Instance**
    - Version: v22.1.8
    - Installation: [CockroachDB Official Website](https://www.cockroachlabs.com/docs/v22.1/install-cockroachdb.html)

2. **Postgres**
    - Installation: [Postgres Official Website](https://www.postgresql.org/download/)

3. **Redis**
    - Installation: [Redis Official Website](https://redis.io/download)

4. **Envfile (Configuration File)**
    - Setup: Create an environment file with required configuration values. You can use the provided template or customize it according to your needs.

5. **Postgres Database**
    - Ensure you have an existing Postgres database with appropriate access privileges for the user mentioned in the environment file.

## CockroachDB Setup

1. To start a single node of CockroachDB in insecure mode, run the following command:
    `cockroach start --insecure --listen-addr=localhost:26257 --http-addr=localhost:8082`
    Once it starts running, you can access the CockroachDB Admin UI by visiting [localhost:8082](http://localhost:8082). Feel free to use any available port of your choice for `--http-addr`.

## Postgres Setup

1. Provide the required environment values in the `.env` file for Postgres. Update the following placeholders with your actual values:
   `DB_TYPE=postgres`
  ` DB_CONNECTION_STRING='host=yourlocalhost user=your_postgres_username password=your_postgres_password dbname=your_db_name port=5432 sslmode=disable'
`
Replace `yourlocalhost`, `your_postgres_username`, `your_postgres_password`, `your_db_name` with your actual Postgres database credentials.

## Redis Setup

1. Setup and run Redis on your computer.

2. Provide the required environment values in the `.env` file for Redis. Update the following placeholders with your actual values:
`REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=`

Replace `localhost` with the appropriate Redis host address and customize `REDIS_PORT` if you are using a different port.

## Preparing Postgres Database

Before your application can work correctly, you need to populate the following tables in your Postgres database:

- asset_classes
- curated_assets
- default_assets
- permissions
- security_questions

Please ensure that you have the necessary data in these tables before running the application.


### You can the go ahead to clone Payment History Engine
   ```bash
   git clone git@github.com:trovotech-technologies/trovo-wallet-payment-history-engine.git
