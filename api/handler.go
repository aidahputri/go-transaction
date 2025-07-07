package api

import (
	// "context"
	"encoding/json"
	"net/http"
	"regexp"

	// "strconv"

	"github.com/aidahputri/go-transaction/model"
	"github.com/aidahputri/go-transaction/repo"
)

var accountNumberRegex = regexp.MustCompile(`^\d{10}$`)

type Handler struct {
	AccountRepo     *repo.Account
	TransactionRepo *repo.Transaction
}

func NewHandler(accountRepo *repo.Account, transactionRepo *repo.Transaction) *Handler {
	return &Handler{
		AccountRepo:     accountRepo,
		TransactionRepo: transactionRepo,
	}
}

// --------- HANDLER ACCOUNT -----------

func (h *Handler) GetAccount(w http.ResponseWriter, r *http.Request) {
	accountNumber := r.URL.Query().Get("accountNumber")
	if accountNumber == "" {
		http.Error(w, "accountNumber is required", http.StatusBadRequest)
		return
	}

	account, err := h.AccountRepo.Get(r.Context(), accountNumber)
	if err != nil {
		http.Error(w, "account not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(account)
}

func (h *Handler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var acc model.Account
	if err := json.NewDecoder(r.Body).Decode(&acc); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if !accountNumberRegex.MatchString(acc.AccountNumber) {
		http.Error(w, "invalid account number format: must be 10 digits", http.StatusBadRequest)
		return
	}
	
	if err := h.AccountRepo.Create(r.Context(), acc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "account created"})
}

func (h *Handler) TopUp(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AccountNumber string  `json:"accountNumber"`
		Amount        float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Amount <= 0 {
		http.Error(w, "amount must be greater than zero", http.StatusBadRequest)
		return
	}

	account, err := h.AccountRepo.Get(r.Context(), req.AccountNumber)
	if err != nil {
		http.Error(w, "account not found", http.StatusNotFound)
		return
	}

	account.Balance += req.Amount
	_, err = h.AccountRepo.Update(r.Context(), account)
	if err != nil {
		http.Error(w, "failed to update balance", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(account)
}

func (h *Handler) BlacklistAccount(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AccountNumber string `json:"accountNumber"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	account, err := h.AccountRepo.Get(r.Context(), req.AccountNumber)
	if err != nil {
		http.Error(w, "account not found", http.StatusNotFound)
		return
	}

	account.Blacklisted = true
	_, err = h.AccountRepo.Update(r.Context(), account)
	if err != nil {
		http.Error(w, "failed to blacklist account", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "account blacklisted"})
}

// --------- HANDLER TRANSACTION -----------

func (h *Handler) Transfer(w http.ResponseWriter, r *http.Request) {
	var tx model.Transaction
	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	fromAcc, err := h.AccountRepo.Get(ctx, tx.FromAccount)
	if err != nil {
		http.Error(w, "sender account not found", http.StatusNotFound)
		return
	}
	toAcc, err := h.AccountRepo.Get(ctx, tx.ToAccount)
	if err != nil {
		http.Error(w, "receiver account not found", http.StatusNotFound)
		return
	}

	if fromAcc.Balance < tx.Amount {
		http.Error(w, "insufficient balance", http.StatusBadRequest)
		return
	}

	// set flag
	if fromAcc.Blacklisted {
		fromAcc.Underwatch = true
	}
	if toAcc.Blacklisted {
		toAcc.Underwatch = true
	}

	fromAcc.Balance -= tx.Amount
	toAcc.Balance += tx.Amount

	// update both accounts
	if _, err = h.AccountRepo.Update(ctx, fromAcc); err != nil {
		http.Error(w, "failed to update sender", http.StatusInternalServerError)
		return
	}
	if _, err = h.AccountRepo.Update(ctx, toAcc); err != nil {
		http.Error(w, "failed to update receiver", http.StatusInternalServerError)
		return
	}

	// insert transaction
	if err := h.TransactionRepo.Create(ctx, tx); err != nil {
		http.Error(w, "failed to record transaction", http.StatusInternalServerError)
		return
	}

	// TODO: Publish to Kafka here if needed

	json.NewEncoder(w).Encode(map[string]string{"message": "transfer successful"})
}
