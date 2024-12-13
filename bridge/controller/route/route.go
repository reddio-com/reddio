package route

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/reddio-com/reddio/bridge/controller/api"
)

// Route routes the APIs
func Route(router *gin.Engine) {
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r := router.Group("bridge/")
	// bridgeApi := api.NewBridgeAPI(db)
	// r.GET("/test", bridgeApi.GetStatus)
	r.POST("/withdrawals", api.L2UnclaimedWithdrawalsByAddressCtl.GetL2UnclaimedWithdrawalsByAddress)
	r.POST("/txsbyaddress", api.TxsByAddressCtl.GetTxsByAddress)

}
