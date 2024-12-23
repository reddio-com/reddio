package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/reddio-com/reddio/bridge/logic"
	"github.com/reddio-com/reddio/bridge/types"
	"gorm.io/gorm"
)

type BridgeAPI struct {
	historyLogic *logic.HistoryLogic
}

func NewBridgeAPI(db *gorm.DB) *BridgeAPI {
	return &BridgeAPI{
		historyLogic: logic.NewHistoryLogic(db),
	}
}

func (s *BridgeAPI) GetStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "Bridge service is running"})
}

// GetL2UnclaimedWithdrawalsByAddress defines the http get method behavior
func (s *BridgeAPI) GetL2UnclaimedWithdrawalsByAddress(ctx *gin.Context) {
	var req types.QueryByAddressRequest
	if err := ctx.ShouldBind(&req); err != nil {
		types.RenderFailure(ctx, types.ErrParameterInvalidNo, err)
		return
	}

	pagedTxs, total, err := s.historyLogic.GetL2UnclaimedWithdrawalsByAddress(ctx, req.Address, req.Page, req.PageSize)
	if err != nil {
		types.RenderFailure(ctx, types.ErrGetL2ClaimableWithdrawalsError, err)
		return
	}

	resultData := &types.ResultData{Results: pagedTxs, Total: total}
	types.RenderSuccess(ctx, resultData)
}
