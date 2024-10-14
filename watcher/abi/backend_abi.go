package backendabi

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var (
	IL1ParentBridgeCoreFacetABI *abi.ABI
	IL2ChildBridgeCoreFacetABI  *abi.ABI

	L1DownwardMessageEventSig common.Hash
	L2UpwardMessageEventSig   common.Hash
)

func init() {
	IL1ParentBridgeCoreFacetABI, _ = IL1ParentBridgeCoreFacetMetaData.GetAbi()
	L1DownwardMessageEventSig = IL1ParentBridgeCoreFacetABI.Events["DownwardMessage"].ID

	IL2ChildBridgeCoreFacetABI, _ = IL2ChildBridgeCoreFacetMetaData.GetAbi()
	L2UpwardMessageEventSig = IL2ChildBridgeCoreFacetABI.Events["UpwardMessage"].ID
}

var IL1ParentBridgeCoreFacetMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"sequence\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"payloadType\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"DownwardMessage\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"payloadType\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"sendDownwardMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}
var IL2ChildBridgeCoreFacetMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"sequence\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"payloadType\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"UpwardMessage\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"payloadType\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"sendUpwardMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

type ChildBridgeCoreFacetUpwardMessageEvent struct {
	Sequence    *big.Int
	PayloadType uint32
	Payload     []byte
	Raw         types.Log // Blockchain specific contextual infos
}
