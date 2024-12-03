package relayer

import (
	"github.com/reddio-com/reddio/bridge/contract"
	"github.com/reddio-com/reddio/bridge/orm"
)

type L2ToL1RelayerInterface interface {
	HandleUpwardMessage([]*orm.CrossMessage) error
}

type L1ToL2RelayerInterface interface {
	HandleDownwardMessageWithSystemCall(msg *contract.ParentBridgeCoreFacetDownwardMessage) error
}
