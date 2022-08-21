# wallet

Wallet is an extrimely simple payment system. This is a test that one company gave me to implement.

The task is to develop a simple API with several features:
* API should allow to add a new wallet. Each wallet has a unique name.
* Clients (wallets in this case) allowed to make transactions to wallet (deposit) from wallet (withdraw) and between wallets (transfer).
* All activity of a wallet should be saved and available for review.
* It should be possible to fetch transactions with some filtering.
* Need an endpoint to download filtered transactions report in CSV.


Psql was chosen as a db for this simlified API.

## Run API

In order to run api for the first time use:

```
make run
```

After downloading psql image and running couple containers migrations has to run.

[golang-migrate](https://github.com/golang-migrate/migrate) tool is being used for db migrations, it should be [installed](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) in the local env.

#### Running migrations

```
make db-migrate
```

## Generate some fake data

There is an endpoint that allows to generate some amount of wallets and transactions:
```
curl -X POST 'http://localhost:8080/api/v1/generate_fake_data?records=500'
```
Use `records` query param to decide how many wallets you need. This request should only be used when there is an empty DB with migrations run.

## Available endpoints:
Can be found in internal/adapters/api/{model_name}/openapi.yaml

File (csv) downloading can be done: 
```
curl -X POST 'http://localhost:8080/api/v1/transactions-report?limit=100&offset=0' \
--header 'Content-Type: text/csv' \
--data-raw '{
    "sender_ids": [
        501
    ]
}' > ~/ex.csv
```