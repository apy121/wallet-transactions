package routes

import (
	"slice/main/controllers"
	"slice/main/types"

	"github.com/gin-gonic/gin"
)

func SetupRouter(service types.WalletService) *gin.Engine {
	r := gin.Default()
	ctrl := controllers.NewWalletController(service)

	r.POST("/v1/wallets", ctrl.CreateWallet)
	r.GET("/v1/wallets", ctrl.GetWalletBalance)
	r.POST("/v1/wallets/add", ctrl.AddMoney)
	r.POST("/v1/wallets/withdraw", ctrl.WithdrawMoney)
	r.POST("/v1/transactions", ctrl.TransferMoney)
	r.GET("/v1/transaction/wallet", ctrl.GetTransactionsForWallet)
	r.GET("/v1/transaction", ctrl.GetTransactionsForUser)

	return r
}
