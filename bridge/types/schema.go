package types

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	Sepolia    = 11155111
	Ethereum   = 1
	ReddioTest = 50341
	Reddio     = 50342

	// Success indicates that the operation was successful.
	Success = 0
	// InternalServerError represents a fatal error occurring on the server.
	InternalServerError = 500
	// ErrParameterInvalidNo represents an error when the parameters are invalid.
	ErrParameterInvalidNo = 40001
	// ErrGetL2ClaimableWithdrawalsError represents an error when trying to get L2 claimable withdrawal transactions.
	ErrGetL2ClaimableWithdrawalsError = 40002
	// ErrGetTxsError represents an error when trying to get transactions by address.
	ErrGetTxsError = 40003
)

type CheckStatus int

const (
	CheckStatusUnChecked CheckStatus = iota
	CheckStatusCheckedStep1
	CheckStatusCheckedStep2
	//99 means the check no need to do
)

type ProcessStatus int

const (
	UnProcessed ProcessStatus = iota + 1
	Processed
	ProcessFailed
)

type EventType int

const (
	QueueTransaction EventType = iota + 1 // 1. QueueTransaction (L1DepositMsgSent)
	L2RelayedMessage                      // 2. L2RelayedMessage (DepositMsgConsumed)
	SentMessage                           // 3. SentMessage (L2withdrawMsgSent)
	L1RelayedMessage                      // 4. L1RelayedMessage (withdrawMsgConsumed)
)

type MessagePayloadType int

const (
	PayloadTypeETH MessagePayloadType = iota
	PayloadTypeERC20
	PayloadTypeERC721
	PayloadTypeERC1155
	PayloadTypeRED
)

type TokenType int

const (
	ETH TokenType = iota
	ERC20
	ERC721
	ERC1155
	RED
)

type TxType int

const (
	TxTypeUnknown TxType = iota
	TxTypeDeposit
	TxTypeWithdraw
	TxTypeRefund
)

type TxStatusType int

// Constants for TxStatusType.
const (
	TxStatusTypeSent TxStatusType = iota
	TxStatusTypeConsumed
	TxStatusTypeDropped
	TxStatusTypeReadyForConsumption
)

// MessageType represents the type of message.
type MessageType int

// Constants for MessageType.
const (
	MessageTypeUnknown MessageType = iota
	MessageTypeL1SentMessage
	MessageTypeL2SentMessage
)

// QueryByAddressRequest the request parameter of address api
type QueryByAddressRequest struct {
	Address  string `json:"address" binding:"required"`
	Page     uint64 `json:"page" binding:"required,min=1"`
	PageSize uint64 `json:"page_size" binding:"required,min=1,max=100"`
}

// ResultData contains return txs and total
type ResultData struct {
	Results []*TxHistoryInfo `json:"results"`
	Total   uint64           `json:"total"`
}

// Response the response schema
type Response struct {
	ErrCode int         `json:"errcode"`
	ErrMsg  string      `json:"errmsg"`
	Data    interface{} `json:"data"`
}

type Message struct {
	PayloadType uint32 `json:"payload_type"` // 0: ETH, 1: ERC20, 2: ERC721, 3: ERC1155, 4: RED
	Payload     string `json:"payload"`
	Nonce       string `json:"nonce"`
}

type ClaimInfo struct {
	From    string         `json:"from"`
	To      string         `json:"to"`
	Value   string         `json:"value"`
	Message Message        `json:"message"`
	Proof   L2MessageProof `json:"proof"`
}

// L2MessageProof is the schema of L2 message proof
type L2MessageProof struct {
	MultiSignProof string `json:"multisign_proof"`
}

// TxHistoryInfo the schema of tx history infos
type TxHistoryInfo struct {
	Hash           string       `json:"hash"`
	MessageHash    string       `json:"message_hash"`
	TokenType      TokenType    `json:"token_type"`    // 0: ETH, 1: ERC20, 2: ERC721, 3: ERC1155, 4: RED
	TokenIDs       []string     `json:"token_ids"`     // only for erc721 and erc1155
	TokenAmounts   []string     `json:"token_amounts"` // for eth and erc20, the length is 1, for erc721 and erc1155, the length could be > 1
	MessageType    MessageType  `json:"message_type"`  // 0: unknown, 1: layer 1 message, 2: layer 2 message
	TxStatus       TxStatusType `json:"tx_status"`
	L1TokenAddress string       `json:"l1_token_address"`
	L2TokenAddress string       `json:"l2_token_address"`
	BlockNumber    uint64       `json:"block_number"`
	ClaimInfo      *ClaimInfo   `json:"claim_info"`
	TxType         TxType       `json:"tx_type"` // 0: unknown, 1: deposit, 2: withdraw, 3: refund
	BlockTimestamp uint64       `json:"block_timestamp"`
}

// func getTxHistoryInfoFromCrossMessage(message *orm.CrossMessage) *types.TxHistoryInfo {
// 	txHistory := &types.TxHistoryInfo{
// 		MessageHash:    message.MessageHash,
// 		TokenType:      btypes.TokenType(message.TokenType),
// 		TokenIDs:       utils.ConvertStringToStringArray(message.TokenIDs),
// 		TokenAmounts:   utils.ConvertStringToStringArray(message.TokenAmounts),
// 		L1TokenAddress: message.L1TokenAddress,
// 		L2TokenAddress: message.L2TokenAddress,
// 		MessageType:    btypes.MessageType(message.MessageType),
// 		TxStatus:       btypes.TxStatusType(message.TxStatus),
// 		BlockTimestamp: message.BlockTimestamp,
// 	}
// 	if txHistory.MessageType == btypes.MessageTypeL1SentMessage {
// 		txHistory.Hash = message.L1TxHash
// 		txHistory.ReplayTxHash = message.L1ReplayTxHash
// 		txHistory.RefundTxHash = message.L1RefundTxHash
// 		txHistory.BlockNumber = message.L1BlockNumber
// 		txHistory.CounterpartChainTx = &types.CounterpartChainTx{
// 			Hash:        message.L2TxHash,
// 			BlockNumber: message.L2BlockNumber,
// 		}
// 	} else {
// 		txHistory.Hash = message.L2TxHash
// 		txHistory.BlockNumber = message.L2BlockNumber
// 		txHistory.CounterpartChainTx = &types.CounterpartChainTx{
// 			Hash:        message.L1TxHash,
// 			BlockNumber: message.L1BlockNumber,
// 		}
// 		if btypes.RollupStatusType(message.RollupStatus) == btypes.RollupStatusTypeFinalized {
// 			txHistory.ClaimInfo = &types.ClaimInfo{
// 				From:    message.MessageFrom,
// 				To:      message.MessageTo,
// 				Value:   message.MessageValue,
// 				Nonce:   strconv.FormatUint(message.MessageNonce, 10),
// 				Message: message.MessageData,
// 				Proof: types.L2MessageProof{
// 					BatchIndex:  strconv.FormatUint(message.BatchIndex, 10),
// 					MerkleProof: "0x" + common.Bytes2Hex(message.MerkleProof),
// 				},
// 				Claimable: true,
// 			}
// 		}
// 	}
// 	return txHistory
// }

// RenderJSON renders response with json
func RenderJSON(ctx *gin.Context, errCode int, err error, data interface{}) {
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}
	renderData := Response{
		ErrCode: errCode,
		ErrMsg:  errMsg,
		Data:    data,
	}
	ctx.JSON(http.StatusOK, renderData)
}

// RenderSuccess renders success response with json
func RenderSuccess(ctx *gin.Context, data interface{}) {
	RenderJSON(ctx, Success, nil, data)
}

// RenderFailure renders failure response with json
func RenderFailure(ctx *gin.Context, errCode int, err error) {
	RenderJSON(ctx, errCode, err, nil)
}

// RenderFatal renders fatal response with json
func RenderFatal(ctx *gin.Context, err error) {
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}
	renderData := Response{
		ErrCode: InternalServerError,
		ErrMsg:  errMsg,
		Data:    nil,
	}
	ctx.Set("errcode", InternalServerError)
	ctx.JSON(http.StatusInternalServerError, renderData)
}
