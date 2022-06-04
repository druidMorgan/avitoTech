package controller

import (
	"avitoTechUsBal/service"
	"context"
	"net/http"
)

type TransController struct {
	ctx      context.Context
	services *service.Manager
	// logger
}

// Создание нового контроллера транзакций
func NewTrans(ctx context.Context, services *service.Manager) *TransController {
	return &TransController{
		ctx:      ctx,
		services: services,
	}
}

// Получение списка транзакций
func (ctrl *TransController) GetUserTransactionHandler(w http.ResponseWriter, r *http.Request) {

}
