package service

import (
	"avitoTechUsBal/model"
	"avitoTechUsBal/store"
	"context"
)

type TransWebService struct {
	ctx   context.Context
	store *store.Store
}

func NewTransWebService(ctx context.Context, store *store.Store) *TransWebService {
	return &TransWebService{ctx: ctx, store: store}
}

// Получение полного списка транзакций пользователя
func (trs *TransWebService) GetUserTransactions(userId uint64) (*[]model.Transaction, error) {
	return nil, nil
}

// Получение полного списка транзакций пользователя отсортированного по дате
func (trs *TransWebService) GetUserTransactionsOrderByDate(userId uint64) (*[]model.Transaction, error) {
	return nil, nil
}

// Получение полного списка транзакций пользователя отсортированного по сумме
func (trs *TransWebService) GetUserTransactionsOrderByAmount(userId uint64) (*[]model.Transaction, error) {
	return nil, nil
}

// Получение списка транзакций пользователя в промежутке дат
func (trs *TransWebService) GetUserTransactionsRange(userId uint64, from uint64, to uint64) (*[]model.Transaction, error) {
	return nil, nil
}
