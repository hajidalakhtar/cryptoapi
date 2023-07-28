package domain

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type Wallet struct {
	Address      string         `json:"address"`
	TokenBalance []TokenBalance `json:"token_balance"`
}

type NewWallet struct {
	Note     string `json:"note"`
	Address  string `json:"addr"`
	Mnemonic string `json:"mnemonic"`
}

type TokenBalance struct {
	Name    string     `json:"name"`
	Amount  *big.Float `json:"amount"`
	Symbol  string     `json:"symbol"`
	LogoUrl string     `json:"logo_url"`
	Decimal int        `json:"decimal"`
}

type Token struct {
	Name                 string `json:"name"`
	Symbol               string `json:"symbol"`
	SmartContractAddress string `json:"smart_contract_address"`
	Decimal              string `json:"decimal"`
	LogoUrl              string `json:"logo_url"`
	TokenAbi             string `json:"token_abi"`
}

type Balances struct {
	Name   string     `json:"name"`
	Amount *big.Float `json:"amount"`
}

type WalletUsecase interface {
	GetBalance(context.Context, []string, string) (Wallet, error)
	GenerateNewWallet(context.Context) (NewWallet, error)
	GetBalanceFromMnemonic(context.Context, string) (Wallet, error)
	Transfer(context.Context, string, string, int, int, uint64) (string, error)
}
type RPCWalletRepository interface {
	GetBalance(context.Context, string, string) (*big.Float, error)
	GenerateNewWallet(context.Context) (string, string, error)
	GetBalanceFromMnemonic(context.Context, string) ([]Balances, string, error)
	Transfer(context.Context, string, common.Address, *big.Int, *big.Int, uint64) (string, error)
}

type SupabaseWalletRepository interface {
	GetToken(context.Context, []string) ([]Token, error)
	AddToken(context.Context, []string) ([]Token, error)
}
