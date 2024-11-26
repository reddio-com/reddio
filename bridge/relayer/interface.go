package relayer

import (
	"github.com/reddio-com/reddio/bridge/contract"
)

type L2ToL1RelayerInterface interface {
	HandleUpwardMessage(msg []*contract.ChildBridgeCoreFacetUpwardMessage) error
}

type L1ToL2RelayerInterface interface {
	HandleDownwardMessageWithSystemCall(msg *contract.ParentBridgeCoreFacetDownwardMessage) error
}
