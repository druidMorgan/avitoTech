package service

import (
	"avitoTechUsBal/store"
	"context"
	"fmt"
)

type Manager struct {
	User  UserService
	Trans TransactionService
}

func NewManager(ctx context.Context, store *store.Store) (*Manager, error) {
	if store == nil {
		return nil, fmt.Errorf("No store provided")
	}

	return &Manager{
		User:  NewUserWebService(ctx, store),
		Trans: NewTransWebService(ctx, store),
	}, nil
}
