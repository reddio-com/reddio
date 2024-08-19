package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ethereum/go-ethereum/log"
)

func GetDefaultBlockManager() *BlockManager {
	return &BlockManager{
		hostUrl: "localhost:7999",
	}
}

type BlockManager struct {
	hostUrl string
}

func (bm *BlockManager) StopBlockChain() {
	_, err := http.Get(fmt.Sprintf("http://%s/api/admin/stop", bm.hostUrl))
	if err != nil {
		log.Warn(err.Error())
	}
}

func (bm *BlockManager) GetBlockTxnCountByIndex(index int) (bool, int, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s/api/receipts_count?block_number=%v", bm.hostUrl, index))
	if err != nil {
		return false, 0, err
	}
	defer resp.Body.Close()
	d, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, 0, err
	}
	r := &response{}
	err = json.Unmarshal(d, &r)
	if err != nil {
		return false, 0, err
	}
	if r.ErrMsg == "block not found" {
		return false, 0, nil
	}
	return true, r.Data, nil
}

type response struct {
	Code   int    `json:"code"`
	ErrMsg string `json:"err_msg"`
	Data   int    `json:"data"`
}
