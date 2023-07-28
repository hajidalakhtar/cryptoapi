package usecase

import (
	"context"
	"cryptoapi/domain"
	"log"
	"math/big"
	"strconv"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type WalletUsecase struct {
	rpcWalletRepo  domain.RPCWalletRepository
	supaWalletRepo domain.SupabaseWalletRepository
	contextTimeout time.Duration
}

func NewWalletUsecase(rwr domain.RPCWalletRepository, swr domain.SupabaseWalletRepository, timeout time.Duration) domain.WalletUsecase {
	return &WalletUsecase{
		rpcWalletRepo:  rwr,
		supaWalletRepo: swr,
		contextTimeout: timeout,
	}
}

func (w *WalletUsecase) Transfer(ctx context.Context, mnemonic string, toAddr string, amount int, gasPrice int, gasLimit uint64) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, w.contextTimeout)
	defer cancel()

	toAddrHex := common.HexToAddress(toAddr)
	amountBigInt := big.NewInt(int64(amount))

	txHash, err := w.rpcWalletRepo.Transfer(ctx, mnemonic, toAddrHex, amountBigInt, nil, gasLimit)
	if err != nil {
		return "", err
	}

	return txHash, nil
}

func (w *WalletUsecase) GetBalanceFromMnemonic(ctx context.Context, mnemonic string) (domain.Wallet, error) {
	//ctx, cancel := context.WithTimeout(ctx, w.contextTimeout)
	//defer cancel()
	//
	//bal, addr, err := w.rpcWalletRepo.GetBalanceFromMnemonic(ctx, mnemonic)
	//if err != nil {
	//	return domain.Wallet{}, err
	//}
	//
	//wallet := domain.Wallet{
	//	Address: addr,
	//	Balance: bal,
	//}
	//
	//return wallet, nil

	panic("implement me")
}

func (w *WalletUsecase) GenerateNewWallet(ctx context.Context) (domain.NewWallet, error) {
	ctx, cancel := context.WithTimeout(ctx, w.contextTimeout)
	defer cancel()

	addr, mnemonic, err := w.rpcWalletRepo.GenerateNewWallet(ctx)
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

func (w *WalletUsecase) GetBalance(ctx context.Context, tokenAddresses []string, address string) (domain.Wallet, error) {
	ctx, cancel := context.WithTimeout(ctx, w.contextTimeout)
	defer cancel()

	tokens, err := w.supaWalletRepo.GetToken(ctx, tokenAddresses)
	if err != nil {
		return domain.Wallet{}, err
	}

	balanceChan := make(chan domain.TokenBalance, len(tokenAddresses))
	var wg sync.WaitGroup

	for _, token := range tokens {
		wg.Add(1)
		go func(token domain.Token) {
			defer wg.Done()

			balance, err := w.rpcWalletRepo.GetBalance(ctx, token.SmartContractAddress, address)
			if err != nil {
				log.Printf("Failed to get token balance for token %s: %v", token.Name, err)
				return
			}

			decimal, err := strconv.Atoi(token.Decimal)
			if err != nil {
				log.Printf("Failed to convert decimal for token %s: %v", token.Name, err)
				return
			}

			balanceChan <- domain.TokenBalance{
				Name:    token.Name,
				Amount:  balance,
				Symbol:  token.Symbol,
				LogoUrl: token.LogoUrl,
				Decimal: decimal,
			}
		}(token)
	}

	go func() {
		wg.Wait()
		close(balanceChan)
	}()

	tokenBalance := []domain.TokenBalance{}
	for balance := range balanceChan {
		tokenBalance = append(tokenBalance, balance)
	}

	wallet := domain.Wallet{
		Address:      address,
		TokenBalance: tokenBalance,
	}

	return wallet, nil

}
