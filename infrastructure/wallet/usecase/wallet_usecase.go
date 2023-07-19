package usecase

import (
	"context"
	"cryptoapi/domain"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type WalletUsecase struct {
	walletRepo     domain.WalletRepository
	contextTimeout time.Duration
}

func NewWalletUsecase(wr domain.WalletRepository, timeout time.Duration) domain.WalletUsecase {
	return &WalletUsecase{
		walletRepo:     wr,
		contextTimeout: timeout,
	}
}

func (w *WalletUsecase) Transfer(ctx context.Context, mnemonic string, toAddr string, amount int, gasPrice int, gasLimit uint64) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, w.contextTimeout)
	defer cancel()

	toAddrHex := common.HexToAddress(toAddr)
	amountBigInt := big.NewInt(int64(amount))

	txHash, err := w.walletRepo.Transfer(ctx, mnemonic, toAddrHex, amountBigInt, nil, gasLimit)
	if err != nil {
		return "", err
	}

	return txHash, nil
}

func (w *WalletUsecase) GetBalanceFromMnemonic(ctx context.Context, mnemonic string) (domain.Wallet, error) {
	ctx, cancel := context.WithTimeout(ctx, w.contextTimeout)
	defer cancel()

	bal, addr, err := w.walletRepo.GetBalanceFromMnemonic(ctx, mnemonic)
	if err != nil {
		return domain.Wallet{}, err
	}

	wallet := domain.Wallet{
		Address: addr,
		Balance: bal,
	}

	return wallet, nil
}

func (w *WalletUsecase) GenerateNewWallet(ctx context.Context) (domain.NewWallet, error) {
	ctx, cancel := context.WithTimeout(ctx, w.contextTimeout)
	defer cancel()

	addr, mnemonic, err := w.walletRepo.GenerateNewWallet(ctx)
	if err != nil {
		return domain.NewWallet{}, err
	}

	newWallet := domain.NewWallet{
		Note:     "Please save your mnemonic, if you lose it, you will lose your wallet.",
		Address:  addr,
		Mnemonic: mnemonic,
	}

	return newWallet, nil
}

func (w *WalletUsecase) GetBalance(ctx context.Context, address string) (domain.Wallet, error) {
	ctx, cancel := context.WithTimeout(ctx, w.contextTimeout)
	defer cancel()

	balance, err := w.walletRepo.GetBalance(ctx, address)
	if err != nil {
		return domain.Wallet{}, err
	}

	wallet := domain.Wallet{
		Address: address,
		Balance: balance,
	}

	return wallet, nil
}
