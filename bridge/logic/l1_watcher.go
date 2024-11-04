package logic

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/reddio-com/reddio/evm"
)

type L1WatcherLogic struct {
	cfg                        *evm.GethConfig
	client                     *ethclient.Client
	addressList                []common.Address
	parser                     *L1EventParser
	l1WatcherLogicFetchedTotal *prometheus.CounterVec
}

func NewL1WatcherLogic(cfg *evm.GethConfig, client *ethclient.Client) *L1WatcherLogic {
	contractAddressList := []common.Address{
		common.HexToAddress(cfg.ParentLayerContractAddress),
	}

	f := &L1WatcherLogic{
		cfg:         cfg,
		client:      client,
		addressList: contractAddressList,
		parser:      NewL1EventParser(cfg, client),
	}

	reg := prometheus.DefaultRegisterer
	f.l1WatcherLogicFetchedTotal = promauto.With(reg).NewCounterVec(prometheus.CounterOpts{
		Name: "L2_fetcher_logic_fetched_total",
		Help: "The total number of events or failed txs fetched in L2 fetcher logic.",
	}, []string{"type"})
	return f
}
