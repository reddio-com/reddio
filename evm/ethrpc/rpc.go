package ethrpc

import (
	"context"
	"errors"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/reddio-com/reddio/evm"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/sourcegraph/conc"
	"github.com/yu-org/yu/core/kernel"
	"net"
	"net/http"
	"time"
)

const SolidityTripod = "solidity"

type EthRPC struct {
	chain     *kernel.Kernel
	cfg       *evm.GethConfig
	srv       *http.Server
	rpcServer *rpc.Server
}

func StartupEthRPC(chain *kernel.Kernel, cfg *evm.GethConfig) {
	if cfg.EnableEthRPC {
		rpcSrv, err := NewEthRPC(chain, cfg)
		if err != nil {
			logrus.Fatalf("init EthRPC server failed, %v", err)
		}
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			defer cancel()
			err = rpcSrv.Serve(ctx)
			if err != nil {
				logrus.Errorf("starknetRPC serves failed, %v", err)
			}
		}()
	}
}

func NewEthRPC(chain *kernel.Kernel, cfg *evm.GethConfig) (*EthRPC, error) {
	s := &EthRPC{
		chain:     chain,
		cfg:       cfg,
		rpcServer: rpc.NewServer(),
	}
	logrus.Debug("Start EthRpc at ", net.JoinHostPort(cfg.EthHost, cfg.EthPort))
	backend := &EthAPIBackend{
		allowUnprotectedTxs: true,
		chain:               chain,
		ethChainCfg:         cfg.ChainConfig,
	}
	apis := GetAPIs(backend)
	for _, api := range apis {
		err := s.rpcServer.RegisterName(api.Namespace, api.Service)
		if err != nil {
			return nil, err
		}
	}

	mux := http.NewServeMux()
	mux.Handle("/", s.rpcServer)

	s.srv = &http.Server{
		Addr:        net.JoinHostPort(cfg.EthHost, cfg.EthPort),
		Handler:     cors.Default().Handler(mux),
		ReadTimeout: 30 * time.Second,
	}

	return s, nil
}

func (s *EthRPC) Serve(ctx context.Context) error {
	errCh := make(chan error)
	defer close(errCh)

	var wg conc.WaitGroup
	defer wg.Wait()
	wg.Go(func() {
		if err := s.srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	})

	select {
	case <-ctx.Done():
		return s.srv.Shutdown(context.Background())
	case err := <-errCh:
		return err
	}
}
