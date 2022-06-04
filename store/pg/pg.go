package pg

/*
   Основной файл пакета, в котором реализуется настройка и запуск БД PostgreSQL.
   Настройки подтягиваются из config (ПОКА ХАРДКОД)
*/

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

const Timeout = 5

// DB is a shortcut structure to a Postgres DB
type DB struct {
	*sql.DB
}

func OpenDB() (*DB, error) {
	connStr := "user=postgres password=root dbname=users sslmode=disable" // TODO Брать из конфига
	pgDB, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(Timeout)) // TODO Проверить!
	defer cancel()

	ans, err := pgDB.ExecContext(ctx, "SELECT 1")
	if err != nil {
		return nil, err
	}

	fmt.Println(ans)

	return &DB{pgDB}, nil
}
