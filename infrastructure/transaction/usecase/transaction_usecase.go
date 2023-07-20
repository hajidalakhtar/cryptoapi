package usecase

import (
	"context"
	"cryptoapi/domain"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type TransactionUsecase struct {
	apitr          domain.APITransactionRepository
	rpctr          domain.RPCTransactionRepository
	contextTimeout time.Duration
}

func NewTransactionUsecase(apitr domain.APITransactionRepository, rpctr domain.RPCTransactionRepository, timeout time.Duration) domain.TransactionUsecase {
	return &TransactionUsecase{
		apitr:          apitr,
		rpctr:          rpctr,
		contextTimeout: timeout,
	}
}

func (w *TransactionUsecase) TransactionHistory(ctx context.Context, address string, sort string, page string, limit string) (domain.TransactionResp, error) {

	ctx, cancel := context.WithTimeout(ctx, w.contextTimeout)
	defer cancel()

	transactions, err := w.apitr.TransactionHistory(ctx, address, sort, page, limit)
	if err != nil {
		return domain.TransactionResp{}, err
	}

	transactionsResp := domain.TransactionResp{
		Address: address,
		Txs:     transactions,
	}

	return transactionsResp, nil

}

func (w *TransactionUsecase) Transfer(ctx context.Context, mnemonic string, toAddr string, amount int, gasPrice int, gasLimit uint64) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, w.contextTimeout)
	defer cancel()

	toAddrHex := common.HexToAddress(toAddr)
	amountBigInt := big.NewInt(int64(amount))

	txHash, err := w.rpctr.Transfer(ctx, mnemonic, toAddrHex, amountBigInt, nil, gasLimit)
	if err != nil {
		return "", err
	}

	return txHash, nil
}
