package repository

import (
	"balance_service/pkg/struct4parse"
	"context"
	"errors"
)

var (
	ErrOfInputData = errors.New("неправильные исходные данные для БД")
)

type Repo interface {
	AddIncome(ctx context.Context, income struct4parse.BalanceWithDesc) error
	AddExpense(ctx context.Context, transaction struct4parse.Transaction) error
	GetBalance(ctx context.Context, income *struct4parse.Balance) error
	AddReserve(ctx context.Context, transaction struct4parse.Transaction) error
	GetAllReserved(ctx context.Context, income *[]struct4parse.Reserve) error
	GetAllBalances(ctx context.Context, income *[]struct4parse.Balance) error
	GetAllTransactions(ctx context.Context, income *[]struct4parse.Transaction, by struct4parse.OrderParams) error
	DisReserve(ctx context.Context, expense struct4parse.Transaction) error
	GetReports(ctx context.Context, income *[]struct4parse.Report, timeDur struct4parse.Time4Report) error
}

type Repository struct {
	Repo
}

func NewRepository(db Repo) *Repository {
	return &Repository{
		Repo: db,
	}
}
