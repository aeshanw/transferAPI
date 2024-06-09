# Account-to-Account Transfer API

## Description

Using golang, develop an internal transfers application that facilitates financial transactions between accounts. This application should provide HTTP endpoints for submitting
transaction details and querying account balances.
A postgres database will be used to maintain transaction logs and account states.

## Assumptions

- Consider the currency is the same for all accounts.
- No need to implement authn or authz (i.e No authentication/authorization)
- We aim for high consistency hence some lag in transfers are acceptable to ensure high consistency (i.e CAP theorem - something needs to be sacrificed https://en.wikipedia.org/wiki/CAP_theorem)
- Negative balance is not supported (i.e no bank-like overdrafts)

## Installation

Requirement:
- Docker installed locally --> https://docs.docker.com/desktop/install/mac-install/
- Docker-compose should be able to run-locally (installed along with Docker-engine)

Note: if you wish to run the golang-application code directly (i.e local-compile/run) you'll need `go 1.22.2`

### Database setup
The initdb should be automatically executed by postgres-docker-compose onstartup but incase it fails to initialize you can just copy the `<project-root>/initdb/init.sql` and execute them in your DB client.

### Running unit-tests in docker
```
cd api
docker build -t my-tester-image -f Dockerfile.tester .
docker run --rm my-tester-image
```

## Run
### Using Docker
While in the project-root folder
```
docker-compose up
```
Ensure the database + API logs seem ok before using POSTMAN.

The API will be running on port 3000

You should be able to acces it via POSTMAN

#### Create new account

`POST http://localhost:3000/accounts`

With Payload
```
{
    "account_id": 124,
    "initial_balance": "100.13344"
}
```

Should return 201 response on success

#### Get account details
`GET http://localhost:3000/accounts/124`

#### Transact between 2 accounts
`POST http://localhost:3000/transactions`
With Payload
```
{
    "source_account_id": 124,
    "destination_account_id": 123,
    "amount": "50.12345"
}
```

You should then be able to query the 123 account via
`GET http://localhost:3000/accounts/123`

And verify the balance has been credited accurately (likewise 124 should be debited)
Example:
```
{
    "account_id": 123,
    "balance": "150.35000"
}
```

debiting beyond the source-account's balance will result in a 400 error 
Example:
```
{
    "status": 400,
    "detail": "bad_request",
    "message": "source account has insufficent funds: finalSourceAccountBalance:-49.89345000000001"
}
```

### (Optional) Using local-run

You need to git-clone this folder into your GOPATH e.g `GOPATH/src/aeshanw.com/<this-project-root>` else your go-compiler will not be able to compile or parse the sourcecode.


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
