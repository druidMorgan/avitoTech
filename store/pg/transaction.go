package pg

import "avitoTechUsBal/model"

// Возвращает список транзакций пользователя в заданный период
func (repo *UserPgRepo) GetUserTransactions(userId uint64) (*[]model.Transaction, error) {
	return nil, nil
}

// Возвращает список транзакций пользователя в заданный период
func (repo *UserPgRepo) GetUserTransactionsRange(userId uint64, fromDate int, toDate int, sort string) (*[]model.Transaction, error) {
	return nil, nil
}

// Возвращает список транзакций пользователя в заданный период
func (repo *UserPgRepo) GetUserTransactionsOrderByDate(userId uint64, fromDate int, toDate int, sort string) (*[]model.Transaction, error) {
	return nil, nil
}

// Возвращает список транзакций пользователя в заданный период
func (repo *UserPgRepo) GetUserTransactionsOrderByAmount(userId uint64, fromDate int, toDate int, sort string) (*[]model.Transaction, error) {
	return nil, nil
}
