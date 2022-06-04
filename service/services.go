package service

import "avitoTechUsBal/model"

// Для работы с сервисом пользователей
type UserService interface {
	GetUserBalance(uint64) (*model.User, error)
	TopUpUserBalance(uint64, uint64) (*model.User, error)
	DebitUserBalance(uint64, uint64) (*model.User, error)
	Transfer(uint64, uint64, uint64) (*model.User, error)
}

// Для работы с сервисом транзакций
type TransactionService interface {
	GetUserTransactions(uint64) (*[]model.Transaction, error)
	GetUserTransactionsOrderByDate(userId uint64) (*[]model.Transaction, error)
	GetUserTransactionsOrderByAmount(userId uint64) (*[]model.Transaction, error)
	GetUserTransactionsRange(uint64, uint64, uint64) (*[]model.Transaction, error)
}
