package models

import "time"

// Wallet represents the wallet table in the database
type Wallet struct {
	ID        int
	UserID    int
	CreatedAt time.Time
	UpdatedAt time.Time
	IsDeleted bool
	DeletedAt *time.Time
	Amount    int // In smallest unit (e.g., paise for INR)
	Currency  string
}

type Transaction struct {
	ID                    int
	SourceID              *int64 // Changed to pointer to allow NULL
	ExternalSourceID      *int64
	DestinationID         *int64 // Already a pointer
	ExternalDestinationID *int64
	Type                  string
	CreatedAt             time.Time
	UpdatedAt             time.Time
	IsDeleted             bool
	DeletedAt             *time.Time
	Amount                int
}

// WalletRequest for creating a wallet
type WalletRequest struct {
	UserID int `json:"userId" binding:"required"`
}

// WalletResponse for returning wallet details
type WalletResponse struct {
	WalletID int `json:"wallet_id"`
	Amount   int `json:"amount,omitempty"`
}

// TransactionRequest for adding, withdrawing, or transferring money
type TransactionRequest struct {
	SourceID      *int64 `json:"sourceId,omitempty"`
	DestinationID *int64 `json:"destinationId,omitempty"`
	Amount        int    `json:"amount" binding:"required"`
}

// TransactionResponse for returning transaction details
type TransactionResponse struct {
	TransactionID int `json:"transactionId"`
}
