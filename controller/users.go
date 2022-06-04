package controller

import (
	"avitoTechUsBal/model"
	"avitoTechUsBal/service"
	"avitoTechUsBal/validator"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type UserController struct {
	ctx      context.Context
	services *service.Manager
	// logger
}

// Создание нового контроллера пользователей
func NewUsers(ctx context.Context, services *service.Manager) *UserController {
	return &UserController{
		ctx:      ctx,
		services: services,
	}
}

// Запрос баланса
func (ctrl *UserController) GetUserBalanceHandler(w http.ResponseWriter, r *http.Request) {
	/*
		if r.Method != "GET" {
			fmt.Println("Method not allowed")
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}*/

	status, err := validator.IsRequestValid(r, "GET", "application/json")
	if err != nil {
		http.Error(w, http.StatusText(status), status)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("GetUserBalanceHandler(): body error %v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// парсим пользователя
	balanceInfo := model.BalanceInfo{}
	if err := json.Unmarshal(body, &balanceInfo); err != nil {
		fmt.Printf("Unmarshal error: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// проводим валидацию
	if balanceInfo.UserId == 0 {
		fmt.Println("balanceInfo.UserId == 0")
		http.Error(w, "\"user_id\" must be greater than zero", http.StatusBadRequest)
		return
	}

	user, err := ctrl.services.User.GetUserBalance(balanceInfo.UserId)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	fmt.Println(user)

	jsonUser, err := json.Marshal(user)
	if err != nil {
		fmt.Printf("Marshal error: %v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// вернем пользователя
	fmt.Fprint(w, string(jsonUser))
}

// Пополнение баланса
func (ctrl *UserController) TopUpHandler(w http.ResponseWriter, r *http.Request) {
	/*
		if r.Method != "POST" {
			fmt.Println("Method not allowed")
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	*/

	status, err := validator.IsRequestValid(r, "POST", "application/json")
	if err != nil {
		http.Error(w, http.StatusText(status), status)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("TopUpHandler(): body error %v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// парсим входные параметры
	topUpInfo := model.TopUpInfo{}
	if err := json.Unmarshal(body, &topUpInfo); err != nil {
		fmt.Printf("Unmarshal error: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// проводим валидацию
	if topUpInfo.UserId == 0 {
		fmt.Println("topUpInfo.UserId == 0")
		http.Error(w, "\"user_id\" must be greater than zero", http.StatusBadRequest)
		return
	}

	if topUpInfo.Amount == 0 {
		fmt.Println("topUpInfo.Amount == 0")
		http.Error(w, "\"amount\" must be greater than zero", http.StatusBadRequest)
		return
	}

	user, err := ctrl.services.User.TopUpUserBalance(topUpInfo.UserId, topUpInfo.Amount)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	fmt.Println(user)

	jsonUser, err := json.Marshal(user)
	if err != nil {
		fmt.Printf("Marshal error: %v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// вернем пользователя
	fmt.Fprint(w, string(jsonUser))
}

// Списание средств
func (ctrl *UserController) DebitHandler(w http.ResponseWriter, r *http.Request) {
	/*
		if r.Method != "POST" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}*/

	status, err := validator.IsRequestValid(r, "POST", "application/json")
	if err != nil {
		http.Error(w, http.StatusText(status), status)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("DebitHandler(): body error %v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// парсим входные параметры
	debitInfo := model.DebitInfo{}
	if err := json.Unmarshal(body, &debitInfo); err != nil {
		fmt.Printf("Unmarshal error: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// проводим валидацию
	if debitInfo.UserId == 0 {
		fmt.Println("debitInfo.UserId == 0")
		http.Error(w, "\"user_id\" must be greater than zero", http.StatusBadRequest)
		return
	}

	if debitInfo.Amount == 0 {
		fmt.Println("debitInfo.Amount == 0")
		http.Error(w, "\"amount\" must be greater than zero", http.StatusBadRequest)
		return
	}

	user, err := ctrl.services.User.DebitUserBalance(debitInfo.UserId, debitInfo.Amount)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	fmt.Println(user)

	jsonUser, err := json.Marshal(user)
	if err != nil {
		fmt.Printf("Marshal error: %v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// вернем пользователя
	fmt.Fprint(w, string(jsonUser))
}

// Перевод средств
func (ctrl *UserController) TransferHandler(w http.ResponseWriter, r *http.Request) {
	/*
		if r.Method != "POST" {
			fmt.Println("Method not allowed")
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}*/

	status, err := validator.IsRequestValid(r, "POST", "application/json")
	if err != nil {
		http.Error(w, http.StatusText(status), status)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("TransferHandler(): body error %v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// парсим входные параметры
	transferInfo := model.TransferInfo{}
	if err := json.Unmarshal(body, &transferInfo); err != nil {
		fmt.Printf("Unmarshal error: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// проводим валидацию
	if transferInfo.FromUserId == 0 {
		fmt.Println("topUpInfo.FromUserId == 0")
		http.Error(w, "user_from_id = 0", http.StatusBadRequest)
		return
	}

	if transferInfo.ToUserId == 0 {
		fmt.Println("transferInfo.ToUserId == 0")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if transferInfo.Amount == 0 {
		fmt.Println("transferInfo.Amount == 0")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if transferInfo.FromUserId == transferInfo.ToUserId {
		fmt.Println("transferInfo.FromUserId == transferInfo.ToUserId ")
		http.Error(w, "Transferring funds to yourself is not possible", http.StatusBadRequest)
		return
	}

	userTo, err := ctrl.services.User.Transfer(transferInfo.FromUserId, transferInfo.ToUserId, transferInfo.Amount)

	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	fmt.Println(userTo) // вернем пользователя, на чей счет совершен перевод

	jsonUser, err := json.Marshal(userTo)
	if err != nil {
		fmt.Printf("Marshal error: %v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// вернем пользователя
	fmt.Fprint(w, string(jsonUser))
}
