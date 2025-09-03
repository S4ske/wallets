package handlers

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"
	"wallets/internal/domain"
	"wallets/internal/service"
)

type WalletService interface {
	Withdraw(walletID string, amount int) error
	Deposit(walletID string, amount int) error
	CreateNewWallet(balance int) error
	GetWallet(walletID string) (*domain.Wallet, error)
	GetWallets() ([]*domain.Wallet, error)
}

type OperationRequest struct {
	WalletID      string `json:"walletId" validate:"required,uuid4"`
	OperationType string `json:"operationType" validate:"required,oneof=DEPOSIT WITHDRAW"`
	Amount        int    `json:"amount" validate:"required,gte=0"`
}

type WalletHandler struct {
	svc      WalletService
	validate *validator.Validate
}

func NewWalletHandler(svc WalletService) *WalletHandler {
	return &WalletHandler{svc: svc, validate: validator.New()}
}

func (h *WalletHandler) RegisterRoutes(r chi.Router) {
	r.Post("/wallet", h.UpdateWallet)
	r.Get("/wallets/{id}", h.GetWalletBalance)
	r.Post("/wallets", h.CreateWallet)
	r.Get("/wallets", h.GetWallets)
}

func (h *WalletHandler) UpdateWallet(w http.ResponseWriter, r *http.Request) {
	var opReq OperationRequest

	if err := json.NewDecoder(r.Body).Decode(&opReq); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(opReq); err != nil {
		var errs validator.ValidationErrors
		errors.As(err, &errs)
		http.Error(w, errs.Error(), http.StatusBadRequest)
		return
	}

	var err error

	switch opReq.OperationType {
	case "DEPOSIT":
		err = h.svc.Deposit(opReq.WalletID, opReq.Amount)
	case "WITHDRAW":
		err = h.svc.Withdraw(opReq.WalletID, opReq.Amount)
	}
	if err != nil {
		if errors.Is(err, service.ErrInsufficientBalance) {
			http.Error(w, "insufficient balance", http.StatusBadRequest)
		} else if errors.Is(err, service.ErrInvalidID) {
			http.Error(w, "wallet with this id is not found", http.StatusNotFound)
		} else {
			http.Error(w, "failed", http.StatusBadRequest)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *WalletHandler) GetWalletBalance(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.validate.Var(id, "required,uuid4"); err != nil {
		var errs validator.ValidationErrors
		errors.As(err, &errs)
		http.Error(w, errs.Error(), http.StatusBadRequest)
		return
	}

	wallet, err := h.svc.GetWallet(id)
	if err != nil {
		if errors.Is(err, service.ErrInvalidID) {
			http.Error(w, "wallet with this id is not found", http.StatusNotFound)
		} else {
			http.Error(w, "failed", http.StatusBadRequest)
		}
		return
	}

	w.Write([]byte(strconv.Itoa(wallet.Balance)))
}

type createRequest struct {
	Balance int `json:"balance" validate:"required,gte=0"`
}

func (h *WalletHandler) CreateWallet(w http.ResponseWriter, r *http.Request) {
	var req createRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		var errs validator.ValidationErrors
		errors.As(err, &errs)
		http.Error(w, errs.Error(), http.StatusBadRequest)
		return
	}

	if err := h.svc.CreateNewWallet(req.Balance); err != nil {
		http.Error(w, "failed", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *WalletHandler) GetWallets(w http.ResponseWriter, r *http.Request) {
	wallets, err := h.svc.GetWallets()
	if err != nil {
		http.Error(w, "failed", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(wallets)
}
