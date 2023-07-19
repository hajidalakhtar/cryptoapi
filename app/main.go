package main

import (
	"cryptoapi/helper"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"time"

	walletHttpDelivery "cryptoapi/infrastructure/wallet/delivery/http"
	walletRepo "cryptoapi/infrastructure/wallet/repository/rpc"
	walletUsecase "cryptoapi/infrastructure/wallet/usecase"
	"github.com/gofiber/fiber/v2"
)

func main() {

	timeoutContext := time.Duration(2) * time.Second

	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "Circle",
		BodyLimit:     2097152,
	})

	client, err := ethclient.Dial("https://data-seed-prebsc-1-s2.bnbchain.org:8545")
	//client, err := ethclient.Dial("https://bsc-dataseed2.binance.org")
	defer client.Close()

	if err != nil {
		log.Fatalf("Failed to connect to the Binance Smart Chain client: %v", err)
	}
	defer client.Close()

	// Init Repository
	wr := walletRepo.NewWalletRepository(client)

	// Init Usecase
	wu := walletUsecase.NewWalletUsecase(wr, timeoutContext)

	// Init Delivery
	walletHttpDelivery.NewWalletHandler(app, wu)

	err = app.Listen(":8090")
	if err != nil {
		log.Fatalf("Failed to listen to the port: %v", err)
	}
	helper.PanicIfError(err)
}
