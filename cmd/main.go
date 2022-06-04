package main

import (
	"avitoTechUsBal"
	"avitoTechUsBal/controller"
	"avitoTechUsBal/service"
	"avitoTechUsBal/store"
	"context"
	"fmt"
	"net/http"
)

func main() {

	if err := run(); err != nil {
		fmt.Println("Run error!")
	}
	fmt.Println("avitoTechUsBal started...")
}

func run() error {
	ctx := context.Background()
	// config
	// ...

	// logger
	// ...

	// Init repository store (with PostgreSQL inside)
	store, err := store.New(ctx)
	if err != nil {
		return fmt.Errorf("store.New(): failed")
	}

	// Init service manager
	serviceManager, err := service.NewManager(ctx, store)
	if err != nil {
		return fmt.Errorf("manager.New failed: %v", err)
	}

	// Init controllers
	userController := controller.NewUsers(ctx, serviceManager)
	transController := controller.NewTrans(ctx, serviceManager)

	//var cfg map[string]string
	//cleanenv.ReadConfig("config.yaml", &cfg)
	//cfg := model.Config{"8090"}

	http.HandleFunc("/balance", userController.GetUserBalanceHandler)
	http.HandleFunc("/top-up", userController.TopUpHandler)
	http.HandleFunc("/debit", userController.DebitHandler)
	http.HandleFunc("/transfer", userController.TransferHandler)
	http.HandleFunc("/transactions", transController.GetUserTransactionHandler)

	srv := new(avitoTechUsBal.Server)
	/*
		go func(s *avitoTechUsBal.Server) {
		}(srv)*/
	if err := srv.Run("8090", nil); err != nil {
		fmt.Println("Server run error!")
	}

	return nil
}
