package model

/*
   Этот пакет описывает модель данных, которыми мы оперируем.
   Он универсален и не зависит от применяемой БД, технологии передачи (JSON, Protobuf, etc)
*/
type Config struct {
	Port string `yaml:"port"`
}

type User struct {
	Id      uint64 `json:"-" db:"id"`
	UserId  uint64 `json:"user_id" db:"user_id"`
	Account uint64 `json:"account" db:"account"`
}

type DBUser struct {
	TableName struct{} `pg:"users"`
	Id        uint64   `pg:"id,notnull,pk"`
	FirstName string   `pg:"firstname,notnull"`
	LastName  string   `pg:"lastname,notnull"`
	Account   uint64   `pg:"account,notnull"`
}

// тип принимаемых данных в POST запросе
type BalanceInfo struct {
	UserId uint64 `json:"user_id"`
}

type TopUpInfo struct {
	UserId uint64 `json:"user_id"`
	Amount uint64 `json:"amount"`
}

type DebitInfo struct {
	UserId uint64 `json:"user_id"`
	Amount uint64 `json:"amount"`
}

type TransferInfo struct {
	FromUserId uint64 `json:"user_from_id"`
	ToUserId   uint64 `json:"user_to_id"`
	Amount     uint64 `json:"amount"`
}

func (dbUser *DBUser) ToWeb() *User {
	return &User{
		Id:      dbUser.Id,
		UserId:  dbUser.Id,
		Account: dbUser.Account,
	}
}

// ToDB converts User to DBUser
func (user *User) ToDB() *DBUser {
	return &DBUser{
		Id:      user.Id,
		Account: user.Account,
	}
}
