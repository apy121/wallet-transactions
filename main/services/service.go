package services

import (
	"errors"
	"slice/main/models"
	"slice/main/repositories"
	"slice/main/types"
	"time"
)

type walletService struct {
	repo types.WalletRepository
}

func NewWalletService(repo types.WalletRepository) types.WalletService {
	return &walletService{repo: repo}
}

// CreateWallet validates the request and creates a wallet
func (s *walletService) CreateWallet(req models.WalletRequest) (models.WalletResponse, error) {
	// Validation
	if req.UserID <= 0 {
		return models.WalletResponse{}, errors.New("user ID must be a positive integer")
	}

	walletID, err := s.repo.CreateWallet(req.UserID)
	if err != nil {
		return models.WalletResponse{}, err
	}
	return models.WalletResponse{WalletID: walletID}, nil
}

// GetWalletBalance validates the wallet ID and fetches the balance
func (s *walletService) GetWalletBalance(walletID int) (models.WalletResponse, error) {
	// Validation
	if walletID <= 0 {
		return models.WalletResponse{}, errors.New("wallet ID must be a positive integer")
	}

	wallet, err := s.repo.FindWalletByID(walletID)
	if err != nil {
		return models.WalletResponse{}, err
	}
	if wallet.ID == 0 {
		return models.WalletResponse{}, errors.New("wallet not found")
	}
	return models.WalletResponse{WalletID: wallet.ID, Amount: wallet.Amount}, nil
}

// AddMoney validates the transaction request and adds money to a wallet
func (s *walletService) AddMoney(req models.TransactionRequest) (models.TransactionResponse, error) {
	// Validation
	if req.DestinationID == nil {
		return models.TransactionResponse{}, errors.New("destination wallet ID is required")
	}
	if *req.DestinationID <= 0 {
		return models.TransactionResponse{}, errors.New("destination wallet ID must be a positive integer")
	}
	if req.Amount <= 0 {
		return models.TransactionResponse{}, errors.New("amount must be a positive integer")
	}
	if req.SourceID != nil && *req.SourceID <= 0 {
		return models.TransactionResponse{}, errors.New("external source ID must be a positive integer if provided")
	}

	tx, err := s.repo.BeginTx()
	if err != nil {
		return models.TransactionResponse{}, err
	}
	defer tx.Rollback()

	wallet, err := s.repo.LockWalletByID(tx, int(*req.DestinationID))
	if err == repositories.ErrWalletLocked {
		return models.TransactionResponse{}, errors.New("wallet is currently locked by another transaction")
	}
	if err != nil {
		return models.TransactionResponse{}, err
	}
	if wallet.ID == 0 {
		return models.TransactionResponse{}, errors.New("destination wallet not found")
	}
	if wallet.Amount+req.Amount > repositories.MaxBalanceLimit {
		return models.TransactionResponse{}, errors.New("exceeds maximum balance limit")
	}

	t := models.Transaction{
		ExternalSourceID: req.SourceID,
		DestinationID:    req.DestinationID,
		Type:             "credit",
		CreatedAt:        time.Now(),
		Amount:           req.Amount,
	}
	transactionID, err := s.repo.CreateTransaction(tx, t)
	if err != nil {
		return models.TransactionResponse{}, err
	}
	if err := s.repo.UpdateWalletBalance(int(*req.DestinationID), req.Amount, tx); err != nil {
		return models.TransactionResponse{}, err
	}
	if err := tx.Commit(); err != nil {
		return models.TransactionResponse{}, err
	}
	return models.TransactionResponse{TransactionID: transactionID}, nil
}

// WithdrawMoney validates the transaction request and withdraws money from a wallet
func (s *walletService) WithdrawMoney(req models.TransactionRequest) (models.TransactionResponse, error) {
	// Validation
	if req.SourceID == nil {
		return models.TransactionResponse{}, errors.New("source wallet ID is required")
	}
	if *req.SourceID <= 0 {
		return models.TransactionResponse{}, errors.New("source wallet ID must be a positive integer")
	}
	if req.Amount <= 0 {
		return models.TransactionResponse{}, errors.New("amount must be a positive integer")
	}
	if req.DestinationID != nil && *req.DestinationID <= 0 {
		return models.TransactionResponse{}, errors.New("external destination ID must be a positive integer if provided")
	}

	tx, err := s.repo.BeginTx()
	if err != nil {
		return models.TransactionResponse{}, err
	}
	defer tx.Rollback()

	wallet, err := s.repo.LockWalletByID(tx, int(*req.SourceID))
	if err == repositories.ErrWalletLocked {
		return models.TransactionResponse{}, errors.New("wallet is currently locked by another transaction")
	}
	if err != nil {
		return models.TransactionResponse{}, err
	}
	if wallet.ID == 0 {
		return models.TransactionResponse{}, errors.New("source wallet not found")
	}
	if wallet.Amount < req.Amount {
		return models.TransactionResponse{}, errors.New("insufficient balance")
	}

	t := models.Transaction{
		SourceID:              req.SourceID,
		ExternalDestinationID: req.DestinationID,
		Type:                  "debit",
		CreatedAt:             time.Now(),
		Amount:                req.Amount,
	}
	transactionID, err := s.repo.CreateTransaction(tx, t)
	if err != nil {
		return models.TransactionResponse{}, err
	}
	if err := s.repo.UpdateWalletBalance(int(*req.SourceID), -req.Amount, tx); err != nil {
		return models.TransactionResponse{}, err
	}
	if err := tx.Commit(); err != nil {
		return models.TransactionResponse{}, err
	}
	return models.TransactionResponse{TransactionID: transactionID}, nil
}

// TransferMoney validates the transaction request and transfers money between wallets
func (s *walletService) TransferMoney(req models.TransactionRequest) (models.TransactionResponse, error) {
	// Validation
	if req.SourceID == nil || req.DestinationID == nil {
		return models.TransactionResponse{}, errors.New("source and destination wallet IDs are required")
	}
	if *req.SourceID <= 0 {
		return models.TransactionResponse{}, errors.New("source wallet ID must be a positive integer")
	}
	if *req.DestinationID <= 0 {
		return models.TransactionResponse{}, errors.New("destination wallet ID must be a positive integer")
	}
	if *req.SourceID == *req.DestinationID {
		return models.TransactionResponse{}, errors.New("source and destination wallet IDs must be different")
	}
	if req.Amount <= 0 {
		return models.TransactionResponse{}, errors.New("amount must be a positive integer")
	}

	tx, err := s.repo.BeginTx()
	if err != nil {
		return models.TransactionResponse{}, err
	}
	defer tx.Rollback()

	sourceID := int(*req.SourceID)
	destID := int(*req.DestinationID)
	lowerID, higherID := sourceID, destID
	if sourceID > destID {
		lowerID, higherID = destID, sourceID
	}

	wallet1, err := s.repo.LockWalletByID(tx, lowerID)
	if err == repositories.ErrWalletLocked {
		return models.TransactionResponse{}, errors.New("wallet is currently locked by another transaction")
	}
	if err != nil {
		return models.TransactionResponse{}, err
	}

	wallet2, err := s.repo.LockWalletByID(tx, higherID)
	if err == repositories.ErrWalletLocked {
		return models.TransactionResponse{}, errors.New("wallet is currently locked by another transaction")
	}
	if err != nil {
		return models.TransactionResponse{}, err
	}

	var sourceWallet, destWallet models.Wallet
	if lowerID == sourceID {
		sourceWallet, destWallet = wallet1, wallet2
	} else {
		sourceWallet, destWallet = wallet2, wallet1
	}

	if sourceWallet.ID == 0 || destWallet.ID == 0 {
		return models.TransactionResponse{}, errors.New("wallet not found")
	}
	if sourceWallet.Amount < req.Amount {
		return models.TransactionResponse{}, errors.New("insufficient balance")
	}
	if destWallet.Amount+req.Amount > repositories.MaxBalanceLimit {
		return models.TransactionResponse{}, errors.New("exceeds maximum balance limit")
	}

	t := models.Transaction{
		SourceID:      req.SourceID,
		DestinationID: req.DestinationID,
		Type:          "debit",
		CreatedAt:     time.Now(),
		Amount:        req.Amount,
	}
	transactionID, err := s.repo.CreateTransaction(tx, t)
	if err != nil {
		return models.TransactionResponse{}, err
	}
	if err := s.repo.UpdateWalletBalance(int(*req.SourceID), -req.Amount, tx); err != nil {
		return models.TransactionResponse{}, err
	}
	if err := s.repo.UpdateWalletBalance(int(*req.DestinationID), req.Amount, tx); err != nil {
		return models.TransactionResponse{}, err
	}
	if err := tx.Commit(); err != nil {
		return models.TransactionResponse{}, err
	}
	return models.TransactionResponse{TransactionID: transactionID}, nil
}

// GetTransactionsForWallet validates the wallet ID and fetches transactions
func (s *walletService) GetTransactionsForWallet(walletID int) ([]models.Transaction, error) {
	// Validation
	if walletID <= 0 {
		return nil, errors.New("wallet ID must be a positive integer")
	}

	return s.repo.FindTransactionsByWallet(walletID)
}

// GetTransactionsForUser validates the request and fetches user transactions
func (s *walletService) GetTransactionsForUser(userID int, txType string, startTime, endTime time.Time) ([]models.Transaction, error) {
	// Validation
	if userID <= 0 {
		return nil, errors.New("user ID must be a positive integer")
	}
	if txType != "" && txType != "credit" && txType != "debit" {
		return nil, errors.New("transaction type must be 'credit' or 'debit' if provided")
	}
	if !startTime.IsZero() && !endTime.IsZero() && startTime.After(endTime) {
		return nil, errors.New("start time must be before end time")
	}

	return s.repo.FindTransactionsByUser(userID, txType, startTime, endTime)
}
