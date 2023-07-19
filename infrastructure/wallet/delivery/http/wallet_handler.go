package http

import (
	"cryptoapi/domain"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type WalletHandler struct {
	walletUsecase domain.WalletUsecase
}

func NewWalletHandler(app *fiber.App, wu domain.WalletUsecase) {
	handler := &WalletHandler{walletUsecase: wu}
	app.Get("/bal/:addr", handler.GetBalance)
	app.Get("/generate/wallet", handler.GenerateNewWallet)
	app.Post("/bal", handler.GetBalanceFromMnemonic)
	app.Post("/transfer/:toaddr", handler.Transfer)
}

func (wh *WalletHandler) GetBalanceFromMnemonic(c *fiber.Ctx) error {
	mnemonic := c.FormValue("mnemonic")
	data, err := wh.walletUsecase.GetBalanceFromMnemonic(c.Context(), mnemonic)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(domain.WebResponse{
		Status:  http.StatusOK,
		Data:    data,
		Message: "SUCCESS",
	})
}

func (wh *WalletHandler) GetBalance(c *fiber.Ctx) error {
	addr := c.Params("addr")
	data, err := wh.walletUsecase.GetBalance(c.Context(), addr)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(domain.WebResponse{
		Status:  http.StatusOK,
		Data:    data,
		Message: "SUCCESS",
	})
}

func (wh *WalletHandler) GenerateNewWallet(c *fiber.Ctx) error {
	data, err := wh.walletUsecase.GenerateNewWallet(c.Context())

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(domain.WebResponse{
		Status:  http.StatusOK,
		Data:    data,
		Message: "SUCCESS",
	})
}

func (wh *WalletHandler) Transfer(c *fiber.Ctx) error {
	toAddr := c.Params("toaddr")
	mnemonic := c.FormValue("mnemonic")
	amount, _ := strconv.Atoi(c.FormValue("amount"))
	gasPrice, _ := strconv.Atoi(c.FormValue("gasprice"))
	gasLimit, _ := strconv.ParseUint(c.FormValue("gaslimit"), 10, 64)

	data, err := wh.walletUsecase.Transfer(c.Context(), mnemonic, toAddr, amount, gasPrice, gasLimit)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.WebResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(domain.WebResponse{
		Status:  http.StatusOK,
		Data:    data,
		Message: "SUCCESS",
	})
}
