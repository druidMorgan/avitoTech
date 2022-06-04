package store

import "avitoTechUsBal/model"

/*
   Здесь описыаются интерфейсы хранилища store.
*/

// Для работы с таблицей users
type UserRepo interface {
	GetUserBalance(uint64) (*model.DBUser, error)
	TopUpUserBalance(uint64, uint64) (*model.DBUser, error)
	DebitUserBalance(uint64, uint64) (*model.DBUser, error)
	Transfer(uint64, uint64, uint64) (*model.DBUser, error)
}

// Для работы с таблицей transactions
type TransRepo interface {
	GetUserTransactions(uint64) (*[]model.Transaction, error)
	GetUserTransactionsRange(uint64, uint64, uint64) (*[]model.Transaction, error)
	GetUserTransactionsOrderByDate(userId uint64) (*[]model.Transaction, error)
	GetUserTransactionsOrderByAmount(userId uint64) (*[]model.Transaction, error)
}
