package domain

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type Wallet struct {
	Address string    `json:"addr"`
	Balance []Balance `json:"balance"`
}

type NewWallet struct {
	Note     string `json:"note"`
	Address  string `json:"addr"`
	Mnemonic string `json:"mnemonic"`
}

type Balance struct {
	Name   string     `json:"name"`
	Amount *big.Float `json:"amount"`
}

type WalletUsecase interface {
	GetBalance(context.Context, string) (Wallet, error)
	GenerateNewWallet(context.Context) (NewWallet, error)
	GetBalanceFromMnemonic(context.Context, string) (Wallet, error)
	Transfer(context.Context, string, string, int, int, uint64) (string, error)
}
type WalletRepository interface {
	GetBalance(context.Context, string) ([]Balance, error)
	GenerateNewWallet(context.Context) (string, string, error)
	GetBalanceFromMnemonic(context.Context, string) ([]Balance, string, error)
	Transfer(context.Context, string, common.Address, *big.Int, *big.Int, uint64) (string, error)
}