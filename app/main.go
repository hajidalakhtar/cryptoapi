package main

import (
	"context"
	"cryptoapi/helper"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jackc/pgx/v5"
	"log"
	"os"
	"time"

	walletHttpDelivery "cryptoapi/infrastructure/wallet/delivery/http"
	pgWalletRepo "cryptoapi/infrastructure/wallet/repository/postgresql"
	rpcwalletRepo "cryptoapi/infrastructure/wallet/repository/rpc"
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

	conn, err := pgx.Connect(context.Background(), "postgres://postgres:@127.0.0.1:5432/postgres")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	client, err := ethclient.Dial("https://data-seed-prebsc-2-s1.bnbchain.org:8545")
	defer client.Close()

	bscscanApi := "https://api.bscscan.com/api"

	// Init Repository
	rpcwr := rpcwalletRepo.NewRPCWalletRepository(client)
	pgwr := pgWalletRepo.NewPgWalletRepository(conn)

	apitr := apiTrRepo.NewAPITransactionRepository(bscscanApi)
	rpctr := rpcTrRepo.NewRPCTransactionRepository(client)

	// Init Usecase
	wu := walletUsecase.NewWalletUsecase(rpcwr, pgwr, timeoutContext)
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
