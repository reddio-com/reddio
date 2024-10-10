package relayer

import "github.com/reddio-com/reddio/watcher/contract"

type BridgeRelayerInterface interface {
	HandleDownwardMessageWithSystemCall(msg *contract.ParentBridgeCoreFacetDownwardMessage) error
	HandleUpwardMessage(msg *contract.ChildBridgeCoreFacetUpwardMessage) error
}
