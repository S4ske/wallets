package tests

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/gavv/httpexpect/v2"
	"net/url"
	"testing"
	"wallets/internal/http-server/handlers"
)

const (
	host     = "localhost:8080"
	walletID = "f512717e-1c3f-460e-a3c0-6901c3c2fbce"
)

var (
	u = url.URL{
		Scheme: "http",
		Host:   host,
		Path:   "api/v1",
	}
)

func TestGetWalletBalance_HappyPath(t *testing.T) {
	e := httpexpect.Default(t, u.String())

	e.GET("/wallets/" + walletID).Expect().Status(200).Text().AsNumber().Ge(0)
}

func TestGetWalletBalance_InvalidID(t *testing.T) {
	e := httpexpect.Default(t, u.String())

	e.GET("/wallets/" + gofakeit.UUID()).Expect().Status(404)
}

func TestChangeWalletBalance_HappyPathDeposit(t *testing.T) {
	e := httpexpect.Default(t, u.String())
	e.POST("/wallet").WithJSON(handlers.OperationRequest{
		WalletID:      walletID,
		OperationType: "DEPOSIT",
		Amount:        1000,
	}).
		Expect().Status(200)
}

func TestChangeWalletBalance_HappyPathWithdraw(t *testing.T) {
	e := httpexpect.Default(t, u.String())
	e.POST("/wallet").WithJSON(handlers.OperationRequest{
		WalletID:      walletID,
		OperationType: "WITHDRAW",
		Amount:        500,
	}).
		Expect().Status(200)
}

func TestChangeWalletBalance_NegativeAmountDeposit(t *testing.T) {
	e := httpexpect.Default(t, u.String())
	e.POST("/wallet").WithJSON(handlers.OperationRequest{
		WalletID:      walletID,
		OperationType: "DEPOSIT",
		Amount:        -1000,
	}).
		Expect().Status(400)
}

func TestChangeWalletBalance_NegativeAmountWithdraw(t *testing.T) {
	e := httpexpect.Default(t, u.String())
	e.POST("/wallet").WithJSON(handlers.OperationRequest{
		WalletID:      walletID,
		OperationType: "WITHDRAW",
		Amount:        -1000,
	}).
		Expect().Status(400)
}

func TestChangeWalletBalance_UnknownOperation(t *testing.T) {
	e := httpexpect.Default(t, u.String())
	e.POST("/wallet").WithJSON(handlers.OperationRequest{
		WalletID:      walletID,
		OperationType: "UNKNOWN_OPERATION",
		Amount:        1000,
	}).
		Expect().Status(400)
}

func TestChangeWalletBalance_InvalidId(t *testing.T) {
	e := httpexpect.Default(t, u.String())
	e.POST("/wallet").WithJSON(handlers.OperationRequest{
		WalletID:      gofakeit.UUID(),
		OperationType: "DEPOSIT",
		Amount:        1000,
	}).
		Expect().Status(404)
}
