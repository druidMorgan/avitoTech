package store

/*
    Пакет store это пакет, как бы обощающий работу с хранилищем (БД).
	Это слой абстракции отделяющий конкретную реализацию того, как храняться данные.
	Мы говорим store "сохрани/верни информацию", не задумываясь о том, как именно это будет происходить.
	Такая конструкция дает возможность легко подменять реализации хранилищ, не меняя интерфейсы store.
*/

import (
	"avitoTechUsBal/store/pg"
	"context"
	"errors"
	"log"
	"time"
)

type Store struct {
	Pg *pg.DB

	User        UserRepo
	Transaction TransRepo
}

func New(ctx context.Context) (*Store, error) {
	pgDB, err := pg.OpenDB()

	// connect to Postgres
	if err != nil {
		return nil, errors.New("pgdb.OpenDB failed")
	}

	// Run Postgres migrations
	if pgDB != nil {
		log.Println("Running PostreSQL migrations...")
		// TODO Add migrations here
	}

	var store Store

	if pgDB != nil {
		store.Pg = pgDB
		go store.keepAlivePg()
		store.User = pg.NewUserRepo(pgDB)
	}

	return &store, nil
}

// KeepAlivePollPeriod is a Pg/MySQL keepalive check time period
const KeepAlivePollPeriod = 3

func (store *Store) keepAlivePg() {
	var err error

	for {
		time.Sleep(time.Second * KeepAlivePollPeriod)
		lostConnect := false
		if store.Pg == nil {
			lostConnect = true
		} else if _, err = store.Pg.Exec("SELECT 1"); err != nil {
			lostConnect = true
		}
		if !lostConnect {
			continue
		}
		log.Println("store.KeepAlive(): Lost PostgreSQL connection. Restoring...")
		store.Pg, err = pg.OpenDB()
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("store.KeepAlive(): PostgreSQL connected")
	}
}
