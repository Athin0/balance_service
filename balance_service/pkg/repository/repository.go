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
	GetAllReserved(ctx context.Context, income *[]struct4parse.Transaction) error
	GetAllBalances(ctx context.Context, income *[]struct4parse.Balance) error
}

type Repository struct {
	Repo
}

func NewRepository(db Repo) *Repository {
	return &Repository{
		Repo: db,
	}
}
