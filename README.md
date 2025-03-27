# API Endpoints
Below are the curl commands to test all 7 endpoints, along with descriptions, request formats, and expected responses. All amounts are in paise (e.g., 50000 = 500 INR).

## 1. Create Wallet
**Endpoint:** `POST /v1/wallets`  
**Description:** Creates a new wallet for a user.  
**Request:**
```bash
curl -X POST http://localhost:8080/v1/wallets \
-H "Content-Type: application/json" \
-d '{"userId": 123}'
```

2. Retrieve Wallet Balance
Endpoint: GET /v1/wallets?wallet_id=<wallet_id>

Description: Fetches the balance of a specific wallet.

```
curl -X GET "http://localhost:8080/v1/wallets?wallet_id=1"
```

3. Add Money
Endpoint: POST /v1/wallets/add

Description: Adds money to a wallet from an external source (e.g., bank).

```
curl -X POST http://localhost:8080/v1/wallets/add \
-H "Content-Type: application/json" \
-d '{"sourceId": 1001, "destinationId": 1, "amount": 50000}'

```


4. Withdraw Money
Endpoint: POST /v1/wallets/withdraw

Description: Withdraws money from a wallet to an external destination.


```
curl -X POST http://localhost:8080/v1/wallets/withdraw \
-H "Content-Type: application/json" \
-d '{"sourceId": 1, "destinationId": 1002, "amount": 20000}'
```


5. Transfer Money
Endpoint: POST /v1/transactions

Description: Transfers money between two wallets.

```
curl -X POST http://localhost:8080/v1/transactions \
-H "Content-Type: application/json" \
-d '{"sourceId": 1, "destinationId": 2, "amount": 30000}'

```


6. Get Transactions for a Wallet
Endpoint: GET /v1/transaction/wallet?wallet_id=<wallet_id>

Description: Retrieves all transactions for a specific wallet.

```
curl -X GET "http://localhost:8080/v1/transaction/wallet?wallet_id=1"

```

7. Get Transactions for a User
Endpoint: GET /v1/transaction?user_id=<user_id>&type=<type>&start_time_stamp=<start>&end_time_stamp=<end>

Description: Fetches transaction history for a user with optional filters.

```
curl -X GET "http://localhost:8080/v1/transaction?user_id=123&type=debit&start_time_stamp=2025-03-27T00:00:00Z&end_time_stamp=2025-03-27T23:59:59Z"

```




