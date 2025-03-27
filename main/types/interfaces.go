package types

import (
	"database/sql"
	"slice/main/models"
	"time"
)

// WalletService defines the service layer interface
type WalletService interface {
	CreateWallet(req models.WalletRequest) (models.WalletResponse, error)
	GetWalletBalance(walletID int) (models.WalletResponse, error)
	AddMoney(req models.TransactionRequest) (models.TransactionResponse, error)
	WithdrawMoney(req models.TransactionRequest) (models.TransactionResponse, error)
	TransferMoney(req models.TransactionRequest) (models.TransactionResponse, error)
	GetTransactionsForWallet(walletID int) ([]models.Transaction, error)
	GetTransactionsForUser(userID int, txType string, startTime, endTime time.Time) ([]models.Transaction, error)
}

// WalletRepository defines the repository layer interface
type WalletRepository interface {
	LockWalletByID(tx *sql.Tx, walletID int) (models.Wallet, error)
	CreateWallet(userID int) (int, error)
	FindWalletByID(walletID int) (models.Wallet, error)
	UpdateWalletBalance(walletID, amount int, tx *sql.Tx) error
	CreateTransaction(tx *sql.Tx, transaction models.Transaction) (int, error)
	FindTransactionsByWallet(walletID int) ([]models.Transaction, error)
	FindTransactionsByUser(userID int, txType string, startTime, endTime time.Time) ([]models.Transaction, error)
	BeginTx() (*sql.Tx, error)
}
