package service

import (
	"avitoTechUsBal/model"
	"avitoTechUsBal/store"
	"context"
	"fmt"
)

type UserWebService struct {
	ctx   context.Context
	store *store.Store
}

// Создает новый веб сервис пользователя
func NewUserWebService(ctx context.Context, store *store.Store) *UserWebService {
	return &UserWebService{ctx: ctx, store: store}
}

// Реализации интерфейсных методов

// Обработка запроса на получение баланса пользователя TODO  не доделано
func (svc *UserWebService) GetUserBalance(userId uint64) (*model.User, error) {
	userDB, err := svc.store.User.GetUserBalance(userId)

	if err != nil {
		return nil, fmt.Errorf("svc.store.User.GetUserBalance error: %v", err)
	}

	if userDB == nil {
		return nil, fmt.Errorf("User '%d' not found", userId)
	}

	return userDB.ToWeb(), nil
}

// Обработка запроса на пополнение баланса пользователя
func (svc *UserWebService) TopUpUserBalance(userId uint64, amount uint64) (*model.User, error) {
	userDB, err := svc.store.User.TopUpUserBalance(userId, amount) // пополнение баланса пользователя

	if err != nil {
		return nil, fmt.Errorf("svc.store.User.TopUpUserBalance error: %v", err)
	}

	return userDB.ToWeb(), nil
}

// Обработка запроса на списание средств со счета пользователя TODO  не доделано
func (svc *UserWebService) DebitUserBalance(userId uint64, amount uint64) (*model.User, error) {
	userDB, err := svc.store.User.DebitUserBalance(userId, amount)

	if err != nil {
		return nil, fmt.Errorf("svc.store.User.DebitUserBalance error: %v", err)
	}

	if userDB == nil {
		return nil, fmt.Errorf("User '%d' not found", userId)
	}

	fmt.Println("Debit")
	return userDB.ToWeb(), nil
}

// Обработка запроса на перевод денег между пользователями
func (svc *UserWebService) Transfer(userId, toUserId, amount uint64) (*model.User, error) {
	userDB, err := svc.store.User.Transfer(userId, toUserId, amount)

	if err != nil {
		return nil, fmt.Errorf("svc.store.User.Transfer() error: %v", err)
	}

	if userDB == nil {
		return nil, fmt.Errorf("User '%d' not found", userId)
	}

	return userDB.ToWeb(), nil
}
