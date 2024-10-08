package relayer

import "github.com/reddio-com/reddio/watcher/contract"

type BridgeRelayerInterface interface {
	HandleDownwardMessage(msg *contract.ParentBridgeCoreFacetDownwardMessage) error
	HandleDownwardMessageWithSystemCall(msg *contract.ParentBridgeCoreFacetDownwardMessage) error
}
