# TROVO WALLET API


## Setup Instructions

Install Go on your device, then follow these instructions below:

- `$ git clone https://github.com/trovotech-technologies/trovo-wallet-payment-history-engine.git` the repo.
- Run `$ go test -v ./...` to ensure everything is running smoothly.
- Create a `.env` file. A sample `.env.example` can be found on the root document.
- Run `$ go run main.go` to serve on `0.0.0.0:8080` (for windows "`localhost:8080`")


## Tech Stack Overview

The APIs are written in Go. Notably, we use [Go Gin](https://github.com/gin-gonic/gin) as our web framework.

Logical components are grouped into folders under the /internal folder.

We strive to use  [REST API guideline standards](https://stackoverflow.blog/2020/03/02/best-practices-for-rest-api-design/) as much as we can. 

Trovo Wallet API has been tested to work with the following databases:
- Postgres

## Error Messages

Trovo Wallet API errors are generally structured as follows

|   name   |  type  |                                                                        description                                                                        |
| -------- | ------ | --------------------------------------------------------------------------------------------------------------------------------------------------------- |
| error     | string    | Type of error. |
| message    | string | Human readable description of the error. Localized per query language.                                                                                                                    |
| data   | string | Extra developer data related to the error  |


## i18n

The project at its current stage does not yet support i18n

## Email

For development purposes, use online fake smtp servers like mailtrap.io