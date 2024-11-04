package relayer

import (
	"github.com/reddio-com/reddio/bridge/contract"
)

type L2ToL1RelayerInterface interface {
	HandleUpwardMessage(msg []*contract.ChildBridgeCoreFacetUpwardMessage) error
}

type L1ToL2RelayerInterface interface {
	HandleDownwardMessage(msg *contract.ParentBridgeCoreFacetDownwardMessage) error
	HandleDownwardMessageWithSystemCall(msg *contract.ParentBridgeCoreFacetDownwardMessage) error
}
