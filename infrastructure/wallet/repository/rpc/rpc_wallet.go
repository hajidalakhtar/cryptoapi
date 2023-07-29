package rpc

import (
	"context"
	"cryptoapi/domain"
	"cryptoapi/helper"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"log"
	"math/big"
)

type RPCWalletRepository struct {
	client *ethclient.Client
}

func NewRPCWalletRepository(client *ethclient.Client) domain.RPCWalletRepository {
	return &RPCWalletRepository{client: client}
}

func (r *RPCWalletRepository) GenerateNewWallet(ctx context.Context) (string, string, error) {
	mnemonic, err := hdwallet.NewMnemonic(256)
	if err != nil {
		log.Fatal(err)
	}

	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		log.Fatal(err)
	}

	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		log.Fatal(err)
	}

	return account.Address.Hex(), mnemonic, nil
}

func (r *RPCWalletRepository) GetBalance(ctx context.Context, tokenAddr string, addr string) (*big.Float, error) {

	address := common.HexToAddress(addr)
	if tokenAddr == "0xb8c77482e45f1f44de1745f52c74426c631bdd52" {
		balance, err := r.client.BalanceAt(context.Background(), address, nil)
		if err != nil {
			return nil, err
		}
		nativeTokenBalance := helper.ConvertBalanceToFloat(balance, 18)
		return nativeTokenBalance, nil
	} //bnb

	balance, err := helper.GetTokenBalance(r.client, addr, tokenAddr)
	return balance, err

}

func (r *RPCWalletRepository) Transfer(ctx context.Context, mnemonic string, toAddr common.Address, amount *big.Int, gasPrice *big.Int, gasLimit uint64) (string, error) {
	chainID, err := r.client.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("Failed to get network ID: %v", err)
	}

	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		log.Fatal(err)
	}

	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, true)
	if err != nil {
		log.Fatal(err)
	}

	nonce, err := r.client.NonceAt(context.Background(), account.Address, nil)
	if err != nil {
		log.Fatalf("Failed to get nonce: %v", err)
	}

	gasPrice, err = r.client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Failed to get gas price: %v", err)
	}

	gasLimit = uint64(21000)
	var data []byte

	tx := types.NewTransaction(nonce, toAddr, amount, gasLimit, gasPrice, data)
	signedTx, err := wallet.SignTx(account, tx, chainID)
	if err != nil {
		log.Fatal(err)
	}

	err = r.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("Failed to send the transaction: %v", err)
	}

	return tx.Hash().Hex(), nil
}
