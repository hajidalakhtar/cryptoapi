package usecase

import (
	"context"
	"cryptoapi/domain"
	"cryptoapi/helper"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/rs/zerolog/log"
	"math/big"
	"strconv"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type WalletUsecase struct {
	rpcWalletRepo  domain.RPCWalletRepository
	pgWalletRepo   domain.PostgresqlWalletRepository
	contextTimeout time.Duration
}

func NewWalletUsecase(rwr domain.RPCWalletRepository, swr domain.PostgresqlWalletRepository, timeout time.Duration) domain.WalletUsecase {
	return &WalletUsecase{
		rpcWalletRepo:  rwr,
		pgWalletRepo:   swr,
		contextTimeout: timeout,
	}
}

func (w *WalletUsecase) AddToken(ctx context.Context, tokenAddr string) (domain.Token, error) {

	data, err := helper.GetCoingeckoTokenApiFromSmartContract(tokenAddr)

	if err != nil {
		log.Warn().Err(err).Msg("Error executing the query")
		return domain.Token{}, err
	}

	decimal := strconv.Itoa(data.DetailPlatforms.BinanceSmartChain.DecimalPlace)

	token := domain.Token{
		Name:                 data.Name,
		Symbol:               data.Symbol,
		SmartContractAddress: tokenAddr,
		Decimal:              decimal,
		LogoURL:              data.Image.Small,
		TokenABI:             "tokeAbi",
	}

	tokenResponse, err := w.pgWalletRepo.AddToken(ctx, token)

	return tokenResponse, err
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

func (w *WalletUsecase) GetBalanceFromMnemonic(ctx context.Context, tokenAddresses []string, mnemonic string) (domain.Wallet, error) {
	ctx, cancel := context.WithTimeout(ctx, w.contextTimeout)
	defer cancel()

	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		log.Warn().Err(err).Msg("Error executing the query")
	}

	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		log.Warn().Err(err).Msg("Error parse derivation path")
	}
	address := account.Address.Hex()
	tokens, err := w.pgWalletRepo.GetTokens(ctx, tokenAddresses)
	if err != nil {
		return domain.Wallet{}, err
	}

	balanceChan := make(chan domain.TokenBalance, len(address))
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
				LogoURL: token.LogoURL,
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

	walletResp := domain.Wallet{
		Address:      address,
		TokenBalance: tokenBalance,
	}

	return walletResp, nil
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
	tokens, err := w.pgWalletRepo.GetTokens(ctx, tokenAddresses)
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
				LogoURL: token.LogoURL,
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
