API Endpoints
Below are the curl commands to test all 7 endpoints, along with descriptions, request formats, and expected responses. All amounts are in paise (e.g., 50000 = 500 INR).

1. Create Wallet
Endpoint: POST /v1/wallets
Description: Creates a new wallet for a user.
Request:
bash

Collapse

Wrap

Copy
curl -X POST http://localhost:8080/v1/wallets \
-H "Content-Type: application/json" \
-d '{"userId": 123}'
Success Response:
json

Collapse

Wrap

Copy
{"wallet_id": 1}
Error Response (Invalid User ID):
json

Collapse

Wrap

Copy
{"error": "user ID must be a positive integer"}
2. Retrieve Wallet Balance
Endpoint: GET /v1/wallets?wallet_id=<wallet_id>
Description: Fetches the balance of a specific wallet.
Request:
bash

Collapse

Wrap

Copy
curl -X GET "http://localhost:8080/v1/wallets?wallet_id=1"
Success Response:
json

Collapse

Wrap

Copy
{"wallet_id": 1, "amount": 0}
Error Response (Invalid Wallet ID):
json

Collapse

Wrap

Copy
{"error": "wallet ID must be a positive integer"}
3. Add Money
Endpoint: POST /v1/wallets/add
Description: Adds money to a wallet from an external source (e.g., bank).
Request:
bash

Collapse

Wrap

Copy
curl -X POST http://localhost:8080/v1/wallets/add \
-H "Content-Type: application/json" \
-d '{"sourceId": 1001, "destinationId": 1, "amount": 50000}'
Success Response:
json

Collapse

Wrap

Copy
{"transactionId": 1}
Error Responses:
Invalid Amount: {"error": "amount must be a positive integer"}
Wallet Locked: {"error": "wallet is currently locked by another transaction"}
4. Withdraw Money
Endpoint: POST /v1/wallets/withdraw
Description: Withdraws money from a wallet to an external destination.
Request:
bash

Collapse

Wrap

Copy
curl -X POST http://localhost:8080/v1/wallets/withdraw \
-H "Content-Type: application/json" \
-d '{"sourceId": 1, "destinationId": 1002, "amount": 20000}'
Success Response:
json

Collapse

Wrap

Copy
{"transactionId": 2}
Error Responses:
Insufficient Balance: {"error": "insufficient balance"}
Invalid Source ID: {"error": "source wallet ID must be a positive integer"}
5. Transfer Money
Endpoint: POST /v1/transactions
Description: Transfers money between two wallets.
Request:
bash

Collapse

Wrap

Copy
curl -X POST http://localhost:8080/v1/transactions \
-H "Content-Type: application/json" \
-d '{"sourceId": 1, "destinationId": 2, "amount": 30000}'
Success Response:
json

Collapse

Wrap

Copy
{"transactionId": 3}
Error Responses:
Same Wallet: {"error": "source and destination wallet IDs must be different"}
Exceeds Limit: {"error": "exceeds maximum balance limit"}
6. Get Transactions for a Wallet
Endpoint: GET /v1/transaction/wallet?wallet_id=<wallet_id>
Description: Retrieves all transactions for a specific wallet.
Request:
bash

Collapse

Wrap

Copy
curl -X GET "http://localhost:8080/v1/transaction/wallet?wallet_id=1"
Success Response:
json

Collapse

Wrap

Copy
[
    {
        "id": 1,
        "source_id": null,
        "external_source_id": 1001,
        "destination_id": 1,
        "external_destination_id": null,
        "type": "credit",
        "created_at": "2025-03-27T12:00:00Z",
        "updated_at": "2025-03-27T12:00:00Z",
        "is_deleted": false,
        "deleted_at": null,
        "amount": 50000
    }
]
Error Response (Invalid Wallet ID):
json

Collapse

Wrap

Copy
{"error": "wallet ID must be a positive integer"}
7. Get Transactions for a User
Endpoint: GET /v1/transaction?user_id=<user_id>&type=<type>&start_time_stamp=<start>&end_time_stamp=<end>
Description: Fetches transaction history for a user with optional filters.
Request:
bash

Collapse

Wrap

Copy
curl -X GET "http://localhost:8080/v1/transaction?user_id=123&type=debit&start_time_stamp=2025-03-27T00:00:00Z&end_time_stamp=2025-03-27T23:59:59Z"
Success Response:
json

Collapse

Wrap

Copy
[
    {
        "id": 2,
        "source_id": 1,
        "external_source_id": null,
        "destination_id": null,
        "external_destination_id": 1002,
        "type": "debit",
        "created_at": "2025-03-27T12:01:00Z",
        "updated_at": "2025-03-27T12:01:00Z",
        "is_deleted": false,
        "deleted_at": null,
        "amount": 20000
    }
]
Error Responses:
Invalid Type: {"error": "transaction type must be 'credit' or 'debit' if provided"}
Invalid Time Range: {"error": "start time must be before end time"}
