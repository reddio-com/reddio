package relayer

import (
	backendabi "github.com/reddio-com/reddio/watcher/abi"
	"github.com/reddio-com/reddio/watcher/contract"
)

type L2ToL1RelayerInterface interface {
	HandleUpwardMessage(msg []*backendabi.ChildBridgeCoreFacetUpwardMessageEvent) error
}

type L1ToL2RelayerInterface interface {
	HandleDownwardMessage(msg *contract.ParentBridgeCoreFacetDownwardMessage) error
	HandleDownwardMessageWithSystemCall(msg *contract.ParentBridgeCoreFacetDownwardMessage) error
}
