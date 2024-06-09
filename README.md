# Account-to-Account Transfer API

## Description

Using golang, develop an internal transfers application that facilitates financial transactions between accounts. This application should provide HTTP endpoints for submitting
transaction details and querying account balances.
A postgres database will be used to maintain transaction logs and account states.

## Assumptions

- Consider the currency is the same for all accounts.
- No need to implement authn or authz

## Installation

Assuming:
- Docker installed locally --> https://docs.docker.com/desktop/install/mac-install/
- Docker-compose should be able to run-locally (installed along with Docker-engine)

Note: if you wish to run the golang-application code directly (i.e local-compile/run) you'll need `go 1.22.2`

### Using Docker
While in the project-root folder
```
docker-compose up
```
Ensure the database + API logs seem ok before using POSTMAN.

The API will be running on port 3000

### (Optional) Using local-run
```
cd api
go mod download
make update-vendor
make run
```

#### Test coverage
```
cd api
make test
```



## API Specifications/Requirement

### Account Creation Endpoint (POST)
Implement an endpoint /accounts that accepts JSON-formatted account ID and account initial balance
#### Sample request body
```
{
"account_id": 123,
"initial_balance": "100.23344"
}
```
Expected response is either an error or an empty response, with a suitable http code.
### Account Query Endpoint (GET)
Implement an endpoint accounts/{account_id} that returns the account and its balance of the specified account in a JSON format
Expected response for account ID 123 (if no error)
```
{
"account_id": 123,
"balance": "100.23344"
}
```
### Transaction Submission Endpoint (POST)
Implement an endpoint /transactions that accepts JSON-formatted transaction details, including the source account ID, destination account ID, and transaction amount. The
system should then process these transactions to update the account balances in the database.
#### Sample request body
```
{
"source_account_id": 123,
"destination_account_id": 456,
"amount": "100.12345"
}
```
Expected response is the transaction body, with a suitable http code.

## Architecture

### Services

- All business logic will be organized into packages here
- All domain models will also be maintained here
- DB Access layer utilities will also be kept here

E.g
- AccountService
    - CreateAccount
    - GetAccount
- TransactionService
    - GetTransaction

### Handlers

- All HTTP response-handling & transformation of biz-logic responses to HTTP Errors or statuses will be done in this layer

## Models
Any common shared structs/request models/biz-models that need to be shared cross services/handlers
