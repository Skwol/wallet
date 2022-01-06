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
export POSTGRESQL_URL="postgres://wallet_user:psw@localhost:5432/wallet_db?sslmode=disable"
migrate -database ${POSTGRESQL_URL} -path db/migrations up 
```
#### Removing migrations

```
export POSTGRESQL_URL="postgres://wallet_user:psw@localhost:5432/wallet_db?sslmode=disable"
migrate -database ${POSTGRESQL_URL} -path db/migrations down 
```

## Generate some fake data

There is an endpoint that allows to generate some amount of wallets and transactions:
```
curl -X POST http://localhost:8080/api/v1/generate_fake_data?records=500
```
Use `records` query param to decide how many wallets you need. This request should only be used when there is an empty DB with migrations run.

## Available endpoints:
1. GET /api/v1/wallets/{record_id} get wallet by id, `record_id` is mandatory.
2. GET /api/v1/wallets?limit=10&offset=0 - get all wallets with pagination using limit and offset, `limit` and `offset` are mandatory.
3. GET /api/v1/wallets-with-transactions/{record_id}?limit=10&offset=0 - get wallet by id and all associated transactions. Pagination for transactions using limit and offset. `limit`, `offset` and `record_id` are mandatory.
4. POST /api/v1/wallets - create wallet with json request
```
curl -X POST 'http://localhost:8080/api/v1/wallets' --header 'Content-Type: application/json' --data-raw '{
    "name": "new wallet name two",
    "balance": 500
}'
```
Response:
```
{"id":503,"name":"new wallet name two","balance":500}
```
Create wallet creates a deposit transaction as a side effect if balance > 0.

5. PATCH /api/v1/wallets/{record_id} - update wallet with json request
```
curl -X PATCH 'http://localhost:8080/api/v1/wallets/503' --header 'Content-Type: application/json' --data-raw '{
    "name": "new wallet name two",
    "balance": 600
}'
```
Response:
```
{"id":503,"name":"new wallet name two","balance":600}
```
Update wallet creates a deposit or withdraw transaction as a side effect if balance is being updated.

6. POST /api/v1/transfers - creates transaction between wallets with json request:
```
curl -X POST 'http://localhost:8080/api/v1/transfers' \
--header 'Content-Type: application/json' \
--data-raw '{
    "amount": 100,
    "sender": {
        "id": 503
    },
    "receiver":  {
        "id": 502
    }
}'
```
Response:
```
{
  "id": 505,
  "amount": 100,
  "timestamp": "2022-01-06T16:20:19.5784222Z",
  "sender": {
    "id": 503,
    "balance": 500
  },
  "receiver": {
    "id": 502,
    "balance": 600
  }
}
```

7. GET /api/v1/transactions/{record_id} get transaction by id, `record_id` is mandatory
8. GET /api/v1/transactions?limit=10&offset=0 get all transactions with pagination using limit and offset, `limit` and `offset` are mandatory.
9. POST /api/v1/transactions?limit=100&offset=0 - request transaction with json filter:
```
curl -X POST 'http://localhost:8080/api/v1/transfers' \
--header 'Content-Type: application/json' \
--data-raw '{
    "sender_ids": [
        1,
        2
    ],
    "receiver_ids": [
        1,
        2
    ],
    "types": [
        "deposit"
    ],
    "timestamp": {
        "from": "2022-01-06T16:03:30Z",
        "to": "2022-01-06T16:03:33Z"
    }, 
    "amount": {
        "from": 100,
        "to": 350
    }
}'
```
Response:
```
[
    {
        "id": 1,
        "sender_id": 1,
        "receiver_id": 1,
        "amount": 153.3353,
        "timestamp": "2022-01-06T16:03:30.618775Z",
        "type": "deposit"
    },
    {
        "id": 2,
        "sender_id": 2,
        "receiver_id": 2,
        "amount": 307.3449,
        "timestamp": "2022-01-06T16:03:30.626132Z",
        "type": "deposit"
    }
]
```
10. POST /api/v1/transactions-report?limit=100&offset=0 - request transaction with json filter and receive response as csv.
