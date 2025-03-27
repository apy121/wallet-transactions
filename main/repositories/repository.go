package repositories

import (
	"database/sql"
	"errors"
	"slice/main/models"
	"time"

	"github.com/go-sql-driver/mysql"
)

const (
	MaxBalanceLimit = 20000000 // 2 lakh INR in paise
)

var (
	// ErrWalletLocked is a custom error for when a wallet is locked by another transaction
	ErrWalletLocked = errors.New("wallet is currently locked by another transaction")
)

type walletRepository struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) *walletRepository {
	return &walletRepository{db: db}
}

func (r *walletRepository) CreateWallet(userID int) (int, error) {
	result, err := r.db.Exec("INSERT INTO wallets (user_id, created_at, amount, currency) VALUES (?, ?, 0, 'INR')", userID, time.Now())
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

func (r *walletRepository) FindWalletByID(walletID int) (models.Wallet, error) {
	var w models.Wallet
	var createdAtStr, updatedAtStr string
	var deletedAtStr *string

	err := r.db.QueryRow("SELECT id, user_id, created_at, updated_at, is_deleted, deleted_at, amount, currency FROM wallets WHERE id = ? AND is_deleted = false", walletID).
		Scan(&w.ID, &w.UserID, &createdAtStr, &updatedAtStr, &w.IsDeleted, &deletedAtStr, &w.Amount, &w.Currency)
	if err == sql.ErrNoRows {
		return models.Wallet{}, nil
	}
	if err != nil {
		return models.Wallet{}, err
	}

	w.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
	if err != nil {
		return models.Wallet{}, err
	}
	w.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAtStr)
	if err != nil {
		return models.Wallet{}, err
	}
	if deletedAtStr != nil {
		deletedAt, err := time.Parse("2006-01-02 15:04:05", *deletedAtStr)
		if err != nil {
			return models.Wallet{}, err
		}
		w.DeletedAt = &deletedAt
	}

	return w, nil
}

func (r *walletRepository) LockWalletByID(tx *sql.Tx, walletID int) (models.Wallet, error) {
	var w models.Wallet
	var createdAtStr, updatedAtStr string
	var deletedAtStr *string

	err := tx.QueryRow("SELECT id, user_id, created_at, updated_at, is_deleted, deleted_at, amount, currency FROM wallets WHERE id = ? AND is_deleted = false FOR UPDATE NOWAIT", walletID).
		Scan(&w.ID, &w.UserID, &createdAtStr, &updatedAtStr, &w.IsDeleted, &deletedAtStr, &w.Amount, &w.Currency)
	if err == sql.ErrNoRows {
		return models.Wallet{}, nil
	}
	if err != nil {
		// Check if the error is due to a lock conflict (MySQL error 3572: "Statement aborted because lock(s) could not be acquired immediately")
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 3572 {
			return models.Wallet{}, ErrWalletLocked
		}
		return models.Wallet{}, err
	}

	w.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
	if err != nil {
		return models.Wallet{}, err
	}
	w.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAtStr)
	if err != nil {
		return models.Wallet{}, err
	}
	if deletedAtStr != nil {
		deletedAt, err := time.Parse("2006-01-02 15:04:05", *deletedAtStr)
		if err != nil {
			return models.Wallet{}, err
		}
		w.DeletedAt = &deletedAt
	}

	return w, nil
}

func (r *walletRepository) UpdateWalletBalance(walletID, amount int, tx *sql.Tx) error {
	_, err := tx.Exec("UPDATE wallets SET amount = amount + ?, updated_at = ? WHERE id = ? AND is_deleted = false", amount, time.Now(), walletID)
	return err
}

func (r *walletRepository) CreateTransaction(tx *sql.Tx, t models.Transaction) (int, error) {
	// update in redis lock applied (wallet_id:true)
	// defer redis unlock
	result, err := tx.Exec(
		"INSERT INTO transactions (source_id, external_source_id, destination_id, external_destination_id, type, created_at, amount) VALUES (?, ?, ?, ?, ?, ?, ?)",
		t.SourceID, t.ExternalSourceID, t.DestinationID, t.ExternalDestinationID, t.Type, t.CreatedAt, t.Amount)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

func (r *walletRepository) FindTransactionsByWallet(walletID int) ([]models.Transaction, error) {
	rows, err := r.db.Query(
		"SELECT id, source_id, external_source_id, destination_id, external_destination_id, type, created_at, updated_at, is_deleted, deleted_at, amount "+
			"FROM transactions WHERE (source_id = ? OR destination_id = ?) AND is_deleted = false", walletID, walletID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		var createdAtStr, updatedAtStr string
		var deletedAtStr *string

		err := rows.Scan(&t.ID, &t.SourceID, &t.ExternalSourceID, &t.DestinationID, &t.ExternalDestinationID, &t.Type, &createdAtStr, &updatedAtStr, &t.IsDeleted, &deletedAtStr, &t.Amount)
		if err != nil {
			return nil, err
		}

		t.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			return nil, err
		}
		t.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAtStr)
		if err != nil {
			return nil, err
		}
		if deletedAtStr != nil {
			deletedAt, err := time.Parse("2006-01-02 15:04:05", *deletedAtStr)
			if err != nil {
				return nil, err
			}
			t.DeletedAt = &deletedAt
		}

		transactions = append(transactions, t)
	}
	return transactions, nil
}

func (r *walletRepository) FindTransactionsByUser(userID int, txType string, startTime, endTime time.Time) ([]models.Transaction, error) {
	query := "SELECT t.id, t.source_id, t.external_source_id, t.destination_id, t.external_destination_id, t.type, t.created_at, t.updated_at, t.is_deleted, t.deleted_at, t.amount " +
		"FROM transactions t JOIN wallets w ON (t.source_id = w.id OR t.destination_id = w.id) " +
		"WHERE w.user_id = ? AND t.is_deleted = false"
	args := []interface{}{userID}
	if txType != "" {
		query += " AND t.type = ?"
		args = append(args, txType)
	}
	if !startTime.IsZero() {
		query += " AND t.created_at >= ?"
		args = append(args, startTime)
	}
	if !endTime.IsZero() {
		query += " AND t.created_at <= ?"
		args = append(args, endTime)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		var createdAtStr, updatedAtStr string
		var deletedAtStr *string

		err := rows.Scan(&t.ID, &t.SourceID, &t.ExternalSourceID, &t.DestinationID, &t.ExternalDestinationID, &t.Type, &createdAtStr, &updatedAtStr, &t.IsDeleted, &deletedAtStr, &t.Amount)
		if err != nil {
			return nil, err
		}

		t.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			return nil, err
		}
		t.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAtStr)
		if err != nil {
			return nil, err
		}
		if deletedAtStr != nil {
			deletedAt, err := time.Parse("2006-01-02 15:04:05", *deletedAtStr)
			if err != nil {
				return nil, err
			}
			t.DeletedAt = &deletedAt
		}

		transactions = append(transactions, t)
	}
	return transactions, nil
}

func (r *walletRepository) BeginTx() (*sql.Tx, error) {
	return r.db.Begin()
}
