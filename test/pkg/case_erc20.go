package pkg

import (
	"github.com/ethereum/go-ethereum/common"
	"math/rand"
	"time"
)

type ERC20Wallet struct {
	Address common.Address `json:"address"`
	Balance uint64         `json:"balance"`
}

type CaseERC20Wallet struct {
	*ERC20Wallet
	TokenCount uint64 `json:"tokenCount"`
}

func (w *ERC20Wallet) Copy() *ERC20Wallet {
	return &ERC20Wallet{
		Address: w.Address,
		Balance: w.Balance,
	}
}

func (c *CaseERC20Wallet) Copy() *CaseERC20Wallet {
	return &CaseERC20Wallet{
		ERC20Wallet: c.ERC20Wallet.Copy(),
		TokenCount:  c.TokenCount,
	}
}

type Erc20TransferManager struct {
	ContractAddr common.Address
}

func NewErc20TransferManager(contractAddr common.Address) *Erc20TransferManager {
	return &Erc20TransferManager{
		ContractAddr: contractAddr,
	}
}

func GenerateCaseERC20Wallets(initialTokenCount uint64, wallets []*ERC20Wallet) []*CaseERC20Wallet {
	c := make([]*CaseERC20Wallet, 0)
	for _, w := range wallets {
		c = append(c, &CaseERC20Wallet{
			ERC20Wallet: w,
			TokenCount:  initialTokenCount,
		})
	}
	return c
}

func (m *Erc20TransferManager) GenerateRandomErc20TransferSteps(stepCount int, wallets []*CaseERC20Wallet) *Erc20TransferCase {
	t := &Erc20TransferCase{
		Original: getCopyERC20(wallets),
		Expect:   getCopyERC20(wallets),
	}
	steps := make([]*Erc20Step, 0)
	r := rand.New(rand.NewSource(time.Now().Unix()))
	curTransfer := 1
	for i := 0; i < stepCount; i++ {
		steps = append(steps, generateRandomERC20Step(r, wallets, m.ContractAddr, curTransfer))
		curTransfer++
	}
	t.Steps = steps
	calculateExpectERC20(t)
	return t
}

func (m *Erc20TransferManager) GenerateErc20TransferSteps(wallets []*CaseERC20Wallet) *Erc20TransferCase {
	t := &Erc20TransferCase{
		Original: getCopyERC20(wallets),
		Expect:   getCopyERC20(wallets),
	}
	steps := make([]*Erc20Step, 0)
	curTransfer := 1
	for i := 0; i < len(wallets); i += 2 {
		steps = append(steps, generateERC20Step(wallets[i], wallets[i+1], m.ContractAddr, curTransfer))
		curTransfer++
	}
	t.Steps = steps
	calculateExpectERC20(t)
	return t
}

func (m *Erc20TransferManager) GenerateSameTargetErc20TransferSteps(stepCount int, wallets []*CaseERC20Wallet, target *CaseERC20Wallet) *Erc20TransferCase {
	t := &Erc20TransferCase{
		Original: getCopyERC20(wallets),
		Expect:   getCopyERC20(wallets),
	}
	steps := make([]*Erc20Step, 0)
	cur := 0
	curTransfer := 1
	for i := 0; i < stepCount; i++ {
		from := wallets[cur]
		steps = append(steps, generateERC20TransferStep(from, target, m.ContractAddr, curTransfer))
		cur++
		if cur >= len(wallets) {
			cur = 0
		}
		curTransfer++
	}
	t.Steps = steps
	calculateExpectERC20(t)
	return t
}

func (tc *Erc20TransferCase) Run(m *WalletManager) error {
	nonceMap := make(map[string]uint64)
	for _, step := range tc.Steps {
		fromAddress := step.From.Address.Hex()
		toAddress := step.To.Address.Hex()
		contractAddress := step.ContractAddr.Hex()

		if _, ok := nonceMap[fromAddress]; ok {
			nonceMap[fromAddress]++
		}
		if err := m.TransferERC20(fromAddress, toAddress, contractAddress, step.Count, nonceMap[fromAddress]); err != nil {
			return err
		}
	}
	return nil
}

func (tc *Erc20TransferCase) AssertExpect(m *WalletManager, wallets []*ERC20Wallet) (map[string]*CaseERC20Wallet, bool, error) {
	got := make(map[string]*CaseERC20Wallet)
	for _, w := range wallets {

		addressStr := w.Address.Hex()

		c, err := m.QueryERC20(addressStr, tc.ContractAddr.Hex())
		if err != nil {
			return nil, false, err
		}
		got[addressStr] = &CaseERC20Wallet{
			ERC20Wallet: w,
			TokenCount:  c,
		}
	}
	if len(tc.Expect) != len(got) {
		return got, false, nil
	}
	for key, value := range got {
		e, ok := tc.Expect[key]
		if !ok {
			return got, false, nil
		}
		if e.TokenCount != value.TokenCount {
			return got, false, nil
		}
	}
	return got, true, nil
}

func calculateExpectERC20(tc *Erc20TransferCase) {
	for _, step := range tc.Steps {
		calculateERC20(step, tc.Expect)
	}
}

func calculateERC20(step *Erc20Step, expect map[string]*CaseERC20Wallet) {
	fromAddress := step.From.Address.Hex()
	toAddress := step.To.Address.Hex()

	fromWallet := expect[fromAddress]
	toWallet := expect[toAddress]

	fromWallet.TokenCount = fromWallet.TokenCount - step.Count
	toWallet.TokenCount = toWallet.TokenCount + step.Count
	expect[fromAddress] = fromWallet
	expect[toAddress] = toWallet
}

func generateRandomERC20Step(r *rand.Rand, wallets []*CaseERC20Wallet, contractAddr common.Address, transfer int) *Erc20Step {
	from := r.Intn(len(wallets))
	to := from + 1
	if to >= len(wallets) {
		to = 0
	}
	return &Erc20Step{
		From:         wallets[from].ERC20Wallet,
		To:           wallets[to].ERC20Wallet,
		Count:        uint64(transfer),
		ContractAddr: contractAddr,
	}
}

func generateERC20Step(from, to *CaseERC20Wallet, contractAddr common.Address, transfer int) *Erc20Step {
	return &Erc20Step{
		From:         from.ERC20Wallet,
		To:           to.ERC20Wallet,
		Count:        uint64(transfer),
		ContractAddr: contractAddr,
	}
}

func generateERC20TransferStep(from, to *CaseERC20Wallet, contractAddr common.Address, transferCount int) *Erc20Step {
	return &Erc20Step{
		From:         from.ERC20Wallet,
		To:           to.ERC20Wallet,
		Count:        uint64(transferCount),
		ContractAddr: contractAddr,
	}
}

func getCopyERC20(wallets []*CaseERC20Wallet) map[string]*CaseERC20Wallet {
	m := make(map[string]*CaseERC20Wallet)
	for _, w := range wallets {
		addressStr := w.Address.Hex()
		m[addressStr] = w.Copy()
	}
	return m
}

type Erc20TransferCase struct {
	Steps []*Erc20Step `json:"steps"`
	// address to wallet
	Original     map[string]*CaseERC20Wallet `json:"original"`
	Expect       map[string]*CaseERC20Wallet `json:"expect"`
	ContractAddr common.Address              `json:"contractAddress"`
}

type Erc20Step struct {
	From         *ERC20Wallet   `json:"from"`
	To           *ERC20Wallet   `json:"to"`
	Count        uint64         `json:"count"`
	ContractAddr common.Address `json:"contractAddress"`
}
