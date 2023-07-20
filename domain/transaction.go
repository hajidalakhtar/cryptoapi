package domain

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type Transactions struct {
	BlockNumber       string `json:"blockNumber"`
	TimeStamp         string `json:"timeStamp"`
	Hash              string `json:"hash"`
	Nonce             string `json:"nonce"`
	BlockHash         string `json:"blockHash"`
	TransactionIndex  string `json:"transactionIndex"`
	From              string `json:"from"`
	To                string `json:"to"`
	Value             string `json:"value"`
	Gas               string `json:"gas"`
	GasPrice          string `json:"gasPrice"`
	IsError           string `json:"isError"`
	TxReceiptStatus   string `json:"txreceipt_status"`
	Input             string `json:"input"`
	ContractAddress   string `json:"contractAddress"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	GasUsed           string `json:"gasUsed"`
	Confirmations     string `json:"confirmations"`
	MethodId          string `json:"methodId"`
	FunctionName      string `json:"functionName"`
}

type TransactionResp struct {
	Address string         `json:"address"`
	Txs     []Transactions `json:"txs"`
}

type TransactionUsecase interface {
	TransactionHistory(context.Context, string, string, string, string) (TransactionResp, error)
	Transfer(context.Context, string, string, int, int, uint64) (string, error)
}

type RPCTransactionRepository interface {
	Transfer(context.Context, string, common.Address, *big.Int, *big.Int, uint64) (string, error)
}

type APITransactionRepository interface {
	TransactionHistory(context.Context, string, string, string, string) ([]Transactions, error)
}
