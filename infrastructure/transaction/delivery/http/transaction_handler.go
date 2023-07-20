package http

import (
	"cryptoapi/domain"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type TransactionHandler struct {
	tru domain.TransactionUsecase
}

func NewTransactionHandler(app *fiber.App, tru domain.TransactionUsecase) {
	handler := &TransactionHandler{
		tru: tru,
	}
	app.Post("/transfer/:toaddr", handler.Transfer)
	app.Get("/history/:addr", handler.TransactionHistory)
}

func (th *TransactionHandler) Transfer(c *fiber.Ctx) error {
	toAddr := c.Params("toaddr")
	mnemonic := c.FormValue("mnemonic")
	amount, _ := strconv.Atoi(c.FormValue("amount"))
	gasPrice, _ := strconv.Atoi(c.FormValue("gasprice"))
	gasLimit, _ := strconv.ParseUint(c.FormValue("gaslimit"), 10, 64)

	data, err := th.tru.Transfer(c.Context(), mnemonic, toAddr, amount, gasPrice, gasLimit)

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

func (th *TransactionHandler) TransactionHistory(c *fiber.Ctx) error {

	addr := c.Params("addr")
	page := c.Query("page")
	limit := c.Query("limit")
	sort := c.Query("sort")
	data, err := th.tru.TransactionHistory(c.Context(), addr, sort, page, limit)

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
