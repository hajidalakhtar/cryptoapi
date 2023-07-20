package main

import (
	"cryptoapi/helper"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"time"

	walletHttpDelivery "cryptoapi/infrastructure/wallet/delivery/http"
	walletRepo "cryptoapi/infrastructure/wallet/repository/rpc"
	walletUsecase "cryptoapi/infrastructure/wallet/usecase"

	transactionHttpDelivery "cryptoapi/infrastructure/transaction/delivery/http"
	apiTrRepo "cryptoapi/infrastructure/transaction/repository/api"
	rpcTrRepo "cryptoapi/infrastructure/transaction/repository/rpc"
	transactionUsecase "cryptoapi/infrastructure/transaction/usecase"

	"github.com/gofiber/fiber/v2"
)

func main() {

	timeoutContext := time.Duration(2) * time.Second

	app := fiber.New(fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "Circle",
		BodyLimit:     2097152,
	})

	client, err := ethclient.Dial("https://data-seed-prebsc-1-s2.bnbchain.org:8545")
	bscscanApi := "https://api.bscscan.com/api"
	defer client.Close()

	if err != nil {
		log.Fatalf("Failed to connect to the Binance Smart Chain client: %v", err)
	}
	defer client.Close()

	// Init Repository
	wr := walletRepo.NewWalletRepository(client)
	apitr := apiTrRepo.NewAPITransactionRepository(bscscanApi)
	rpctr := rpcTrRepo.NewRPCTransactionRepository(client)

	// Init Usecase
	wu := walletUsecase.NewWalletUsecase(wr, timeoutContext)
	tru := transactionUsecase.NewTransactionUsecase(apitr, rpctr, timeoutContext)

	// Init Delivery
	walletHttpDelivery.NewWalletHandler(app, wu)
	transactionHttpDelivery.NewTransactionHandler(app, tru)

	err = app.Listen(":8090")
	if err != nil {
		log.Fatalf("Failed to listen to the port: %v", err)
	}
	helper.PanicIfError(err)
}
