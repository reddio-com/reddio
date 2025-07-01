package checker

import (
	"context"
	"fmt"
	"time"

	"github.com/reddio-com/reddio/checker/client"
)

type Checker struct {
	c          *client.NodeClient
	CurrStatus map[string]uint64
	currBlock  uint64
}

func NewChecker(host string) *Checker {
	return &Checker{
		c: client.NewNodeClient(host),
	}
}

func (c *Checker) Run(ctx context.Context) error {
	if err := c.Init(); err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			latestBlock, err := c.c.GetLatestBlock()
			if err != nil {
				return err
			}
			if latestBlock > c.currBlock {
				if err := c.CheckBlock(c.currBlock); err != nil {
					return err
				}
				c.currBlock++
			} else {
				time.Sleep(3 * time.Second)
			}
		}
	}
}

func (c *Checker) Init() error {
	latestBlock, err := c.c.GetLatestBlock()
	if err != nil {
		return err
	}
	c.currBlock = latestBlock
	return nil
}

func (c *Checker) CheckBlock(blockNum uint64) error {
	c.CurrStatus = make(map[string]uint64)
	addrs := make(map[string]struct{})
	txns, err := c.c.GetBlockByNumber(blockNum)
	if err != nil {
		return fmt.Errorf("GetBlockByNumber %v err: %v", blockNum, err)
	}
	for _, txn := range txns {
		if IsTransfer(txn) {
			addrs[txn.From] = struct{}{}
			addrs[txn.To] = struct{}{}
		}
	}
	before, err := c.GetBlockBalance(addrs, blockNum-1)
	if err != nil {
		return err
	}
	cal := make(map[string]uint64)
	for addr, v := range before {
		cal[addr] = v
	}
	now, err := c.GetBlockBalance(addrs, blockNum)
	if err != nil {
		return err
	}
	for _, txn := range txns {
		if IsTransfer(txn) {
			cal[txn.From] -= txn.Value + txn.Gas*txn.GasPrice
			cal[txn.To] += txn.Value
		}
	}
	for addr, v := range cal {
		fmt.Println(fmt.Sprintf("addr: %s, expect: %d, actual: %d", addr, v, now[addr]))
	}
	return nil
}

func (c *Checker) GetBlockBalance(addrs map[string]struct{}, blockNum uint64) (map[string]uint64, error) {
	blockBalance := map[string]uint64{}
	for addr := range addrs {
		v, err := c.c.GetBalanceByBlock(blockNum, addr)
		if err != nil {
			return nil, err
		}
		blockBalance[addr] = v
	}
	return blockBalance, nil
}

func IsTransfer(txn *client.BlockTransaction) bool {
	if txn.Input != "0x" {
		return false
	}
	return true
}
