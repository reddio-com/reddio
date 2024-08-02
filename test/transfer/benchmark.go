package transfer

import (
	"encoding/json"
	"fmt"
	"github.com/reddio-com/reddio/test/pkg"
	"github.com/yu-org/yu/core/kernel"
	"io"
	"net/http"
)

type BenchTestCase struct {
	CaseName string
}

func NewBenchTestCase(name string) *BenchTestCase {
	return &BenchTestCase{CaseName: name}
}

func (b *BenchTestCase) Run(m *pkg.WalletManager) error {
	//TODO: Benchmark transfer here

	// get receipt count
	resp, err := http.Get("http://localhost:7999/api/receipts_count?block_number=3")
	if err != nil {
		return err
	}

	byt, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	apiResp := new(kernel.APIResponse)
	err = json.Unmarshal(byt, apiResp)
	if err != nil {
		return err
	}
	fmt.Println("receipt count = ", apiResp.Data)
	return nil
}

func (b *BenchTestCase) Name() string {
	return b.CaseName
}
