package main

import (
	"cryptoapi/helper"
	"github.com/ethereum/go-ethereum/ethclient"
	supa "github.com/nedpals/supabase-go"

	"log"
	"time"

	walletHttpDelivery "cryptoapi/infrastructure/wallet/delivery/http"
	rpcwalletRepo "cryptoapi/infrastructure/wallet/repository/rpc"
	supawalletRepo "cryptoapi/infrastructure/wallet/repository/supabase"
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

	client, err := ethclient.Dial("https://bsc.publicnode.com")

	defer client.Close()

	supabaseUrl := "https://egkuuffuumguwvalmvdk.supabase.co"
	supabaseKey := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImVna3V1ZmZ1dW1ndXd2YWxtdmRrIiwicm9sZSI6InNlcnZpY2Vfcm9sZSIsImlhdCI6MTY4OTY4OTMwNywiZXhwIjoyMDA1MjY1MzA3fQ.v6P8boW3y00R03uZOmgTOVSN_k931Y8ymJuUcUm9nvo"
	supabase := supa.CreateClient(supabaseUrl, supabaseKey)

	bscscanApi := "https://api.bscscan.com/api"

	// Init Repository
	rpcwr := rpcwalletRepo.NewRPCWalletRepository(client)
	supawr := supawalletRepo.NewSupabaseWalletRepository(supabase)

	apitr := apiTrRepo.NewAPITransactionRepository(bscscanApi)
	rpctr := rpcTrRepo.NewRPCTransactionRepository(client)

	// Init Usecase
	wu := walletUsecase.NewWalletUsecase(rpcwr, supawr, timeoutContext)
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
