package domain

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"math/big"
	"time"
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
	LogoURL string     `json:"logo_url"`
	Decimal int        `json:"decimal"`
}

type Token struct {
	Id                   uuid.UUID `json:"id"`
	Name                 string    `json:"name"`
	Symbol               string    `json:"symbol"`
	SmartContractAddress string    `json:"smart_contract_address"`
	Decimal              string    `json:"decimal"`
	LogoURL              string    `json:"logo_url"`
	TokenABI             string    `json:"token_abi"`
	CreatedAt            time.Time `json:"-"`
	UpdatedAt            time.Time `json:"-"` // Use null.Time for nullable time.Time field.
}

type Balances struct {
	Name   string     `json:"name"`
	Amount *big.Float `json:"amount"`
}

type WalletUsecase interface {
	GetBalance(context.Context, []string, string) (Wallet, error)
	GenerateNewWallet(context.Context) (NewWallet, error)
	GetBalanceFromMnemonic(context.Context, []string, string) (Wallet, error)
	Transfer(context.Context, string, string, int, int, uint64) (string, error)
	AddToken(context.Context, string) (Token, error)
}
type RPCWalletRepository interface {
	GetBalance(context.Context, string, string, string) (*big.Float, error)
	GenerateNewWallet(context.Context) (string, string, error)
	Transfer(context.Context, string, common.Address, *big.Int, *big.Int, uint64) (string, error)
}

type PostgresqlWalletRepository interface {
	GetTokens(context.Context, []string) ([]Token, error)
	AddToken(context.Context, Token) (Token, error)
}
