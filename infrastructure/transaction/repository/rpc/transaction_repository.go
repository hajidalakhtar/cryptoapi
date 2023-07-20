package rpc

import (
	"context"
	"cryptoapi/domain"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"log"
	"math/big"
)

type RPCTransactionRepository struct {
	client *ethclient.Client
}

func NewRPCTransactionRepository(client *ethclient.Client) domain.RPCTransactionRepository {
	return &RPCTransactionRepository{client: client}
}

func (r *RPCTransactionRepository) Transfer(ctx context.Context, mnemonic string, toAddr common.Address, amount *big.Int, gasPrice *big.Int, gasLimit uint64) (string, error) {
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
