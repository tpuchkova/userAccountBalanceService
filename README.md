## Prerequisites
- The **Go** programming language https://golang.org/dl/
- **Docker Desktop** https://www.docker.com/products/docker-desktop
- **GoLand** or similar IDE https://www.jetbrains.com/go/promo/?source=google&medium=cpc&campaign=10156131500&gclid=CjwKCAjwsNiIBhBdEiwAJK4khrn3IDTmD-Xv1BFZ9HQeeSUwIeIFaG69dxoHLW1ACvjxdrZxD5Dn9RoCpXQQAvD_BwE
- **Postman** or similar API client https://www.postman.com/


## Installing&Running
- Clone thist repo https://github.com/tpuchkova/userAccountBalanceService.git
- Start docker desktop app
- Install make `scoop install make`
- Open your terminal and run `make run` command to create and run docker containers
- Run migrations `make migrate`

## Usage
- Example of the POST request

curl --location 'localhost:8080/api/transaction' \
--header 'Source-Type: game' \
  --header 'Content-Type: application/json' \
  --data '{
    "state" : "win",
    "amount": "1",
    "transaction_id" : "93b1a8dd-3ee7-4958-ae83-e3f60ff9129f"
}'

Source-Type header can be game, server or payment.
State can be "win" or "lost".
Win requests increases the user balance. Lost requests decreases user balance.
Each request (with the same transaction id) can be processed only once.

Every 5 minutes 10 latest odd records will be canceled and balance will be corrected by the application. You can see it in application logs.
Canceled records can not be processed twice
