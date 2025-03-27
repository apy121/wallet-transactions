package controllers

import (
	"net/http"
	"slice/main/models"
	"slice/main/types"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type WalletController struct {
	service types.WalletService
}

func NewWalletController(service types.WalletService) *WalletController {
	return &WalletController{service: service}
}

func (ctrl *WalletController) CreateWallet(c *gin.Context) {
	var req models.WalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := ctrl.service.CreateWallet(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (ctrl *WalletController) GetWalletBalance(c *gin.Context) {
	walletID, _ := strconv.Atoi(c.Query("wallet_id"))
	if walletID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wallet_id is required"})
		return
	}
	resp, err := ctrl.service.GetWalletBalance(walletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (ctrl *WalletController) AddMoney(c *gin.Context) {
	var req models.TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := ctrl.service.AddMoney(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (ctrl *WalletController) WithdrawMoney(c *gin.Context) {
	var req models.TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := ctrl.service.WithdrawMoney(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (ctrl *WalletController) TransferMoney(c *gin.Context) {
	var req models.TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := ctrl.service.TransferMoney(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (ctrl *WalletController) GetTransactionsForWallet(c *gin.Context) {
	walletID, _ := strconv.Atoi(c.Query("wallet_id"))
	if walletID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wallet_id is required"})
		return
	}
	transactions, err := ctrl.service.GetTransactionsForWallet(walletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, transactions)
}

func (ctrl *WalletController) GetTransactionsForUser(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))
	txType := c.Query("type")
	startTimeStr := c.Query("start_time_stamp")
	endTimeStr := c.Query("end_time_stamp")

	var startTime, endTime time.Time
	if startTimeStr != "" {
		startTime, _ = time.Parse(time.RFC3339, startTimeStr)
	}
	if endTimeStr != "" {
		endTime, _ = time.Parse(time.RFC3339, endTimeStr)
	}

	transactions, err := ctrl.service.GetTransactionsForUser(userID, txType, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, transactions)
}
