package rpc

import (
	"context"
	"cryptoapi/domain"
	"cryptoapi/infrastructure/transaction/repository/helper"
	"encoding/json"
	"fmt"
	"net/http"
)

type APITransactionRepository struct {
	bscUrl string
}

type Response struct {
	Status  string                `json:"status"`
	Message string                `json:"message"`
	Result  []domain.Transactions `json:"result"`
}

func NewAPITransactionRepository(bscUrl string) domain.APITransactionRepository {
	return &APITransactionRepository{bscUrl: bscUrl}
}

func (A APITransactionRepository) TransactionHistory(ctx context.Context, address string, sort string, page string, offset string) ([]domain.Transactions, error) {
	apiKey := helper.RandomApiKey()
	url := A.bscUrl + "?module=account&action=txlist&address=" + address + "&startblock=0&endblock=99999999&page=" + page + "&offset=" + offset + "&sort=asc&apikey=" + apiKey
	response, err := http.Get(url)

	if err != nil {
		fmt.Println("Error saat melakukan request:", err)
		return nil, err
	}
	defer response.Body.Close()

	var responseData Response
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		fmt.Println("Error saat parsing response:", err)
		return nil, err
	}

	return responseData.Result, nil
}
