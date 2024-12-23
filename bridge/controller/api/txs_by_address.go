package api

import (
	"github.com/reddio-com/reddio/bridge/types"

	"github.com/gin-gonic/gin"
	"github.com/reddio-com/reddio/bridge/logic"
	"gorm.io/gorm"
)

// TxsByAddressController the controller of GetTxsByAddress
type TxsByAddressController struct {
	historyLogic *logic.HistoryLogic
}

// NewTxsByAddressController create new TxsByAddressController
func NewTxsByAddressController(db *gorm.DB) *TxsByAddressController {
	return &TxsByAddressController{
		historyLogic: logic.NewHistoryLogic(db),
	}
}

// GetTxsByAddress defines the http get method behavior
func (c *TxsByAddressController) GetTxsByAddress(ctx *gin.Context) {
	var req types.QueryByAddressRequest
	if err := ctx.ShouldBind(&req); err != nil {
		types.RenderFailure(ctx, types.ErrParameterInvalidNo, err)
		return
	}

	pagedTxs, total, err := c.historyLogic.GetTxsByAddress(ctx, req.Address, req.Page, req.PageSize)
	if err != nil {
		types.RenderFailure(ctx, types.ErrGetTxsError, err)
		return
	}

	resultData := &types.ResultData{Results: pagedTxs, Total: total}
	types.RenderSuccess(ctx, resultData)
}
