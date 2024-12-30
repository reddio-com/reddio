package relayer

import (
	"github.com/reddio-com/reddio/bridge/contract"
	"github.com/reddio-com/reddio/bridge/orm"
)

type L2ToL1RelayerInterface interface {
	HandleUpwardMessage([]*orm.CrossMessage, map[uint64]uint64) error
}

type L1ToL2RelayerInterface interface {
	HandleDownwardMessageWithSystemCall(msg *contract.ParentBridgeCoreFacetQueueTransaction) error
	HandleRelayerMessage(msg *contract.UpwardMessageDispatcherFacetRelayedMessage) error
}
