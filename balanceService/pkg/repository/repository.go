package repository

import (
	"balance_service/pkg/utils"
	"context"
	"errors"
)

var (
	ErrOfInputData = errors.New("неправильные исходные данные для БД")
)

type Repo interface {
	AddIncome(ctx context.Context, income utils.BalanceWithDesc) error
	AddExpense(ctx context.Context, transaction utils.Transaction) error
	GetBalance(ctx context.Context, income *utils.Balance) error
	AddReserve(ctx context.Context, transaction utils.Transaction) error
	GetAllReserved(ctx context.Context, income *[]utils.Reserve) error
	GetAllBalances(ctx context.Context, income *[]utils.Balance) error
	GetAllTransactions(ctx context.Context, income *[]utils.Transaction, by utils.OrderParams) error
	DisReserve(ctx context.Context, expense utils.Transaction) error
	GetReports(ctx context.Context, income *[]utils.Report, timeDur utils.Time4Report) error
}

type Repository struct {
	Repo
}

func NewRepository(db Repo) *Repository {
	return &Repository{
		Repo: db,
	}
}
