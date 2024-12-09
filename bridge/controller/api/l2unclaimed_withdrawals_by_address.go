package api

import (
	"github.com/reddio-com/reddio/bridge/types"

	"github.com/gin-gonic/gin"
	"github.com/reddio-com/reddio/bridge/logic"
	"gorm.io/gorm"
)

// L2UnclaimedWithdrawalsByAddressController the controller of GetL2UnclaimedWithdrawalsByAddress
type L2UnclaimedWithdrawalsByAddressController struct {
	historyLogic *logic.HistoryLogic
}

// NewL2UnclaimedWithdrawalsByAddressController create new L2UnclaimedWithdrawalsByAddressController
func NewL2UnclaimedWithdrawalsByAddressController(db *gorm.DB) *L2UnclaimedWithdrawalsByAddressController {
	return &L2UnclaimedWithdrawalsByAddressController{
		historyLogic: logic.NewHistoryLogic(db),
	}
}

// GetL2UnclaimedWithdrawalsByAddress defines the http get method behavior
func (c *L2UnclaimedWithdrawalsByAddressController) GetL2UnclaimedWithdrawalsByAddress(ctx *gin.Context) {
	//fmt.Println("GetL2UnclaimedWithdrawalsByAddress")
	var req types.QueryByAddressRequest
	if err := ctx.ShouldBind(&req); err != nil {
		types.RenderFailure(ctx, types.ErrParameterInvalidNo, err)
		return
	}

	pagedTxs, total, err := c.historyLogic.GetL2UnclaimedWithdrawalsByAddress(ctx, req.Address, req.Page, req.PageSize)
	if err != nil {
		types.RenderFailure(ctx, types.ErrGetL2ClaimableWithdrawalsError, err)
		return
	}

	resultData := &types.ResultData{Results: pagedTxs, Total: total}
	types.RenderSuccess(ctx, resultData)
}
