package pg

/*
    Это файл реализации конкретный методов работы с БД.
	Здесь хранятся реализации интерфейсов работы с БД:
	 - пополнение, списание, запрос баланса, получение списка транзакций.
	Мы используем модель, которая описывается данные, которыми мы оперируем (User, Transaction).
	Какую бы реализацию БД мы ни выбрали, она (реализация) должна соответствововать модели и разумеется интерфейсу

	Таким образом, пакет pg - это только одна из возможных реализаций работы с хранилищем store.
	То есть store находится выше по иерархии и использует у себя реализацию pg. Если мы захотим сменить хранилище, например,
	с PostgreSQL на MySQL, то нужно будет создать реализацию (пакет), которая обращается к MySQL, используя все те же модели User, Transaction
*/

/*
    Перевод денег реализован на слое БД, потому что используются транзакции.
	Использование транзакций на слое сервиса затруднительно
*/

import (
	"avitoTechUsBal/model"
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

type UserPgRepo struct {
	db *DB
}

func NewUserRepo(db *DB) *UserPgRepo {
	return &UserPgRepo{db: db}
}

// Возвращает баланс пользователя из БД по ключу
func (repo *UserPgRepo) GetUserBalance(userId uint64) (*model.DBUser, error) {

	ctx := context.Background()

	stmt, err := repo.db.PrepareContext(ctx, "SELECT balance FROM users WHERE user_id = $1 FOR UPDATE") // TODO FOR UPDATE CHECK
	if err != nil {
		fmt.Printf("UserPgRepo.GetUserBalance(). Prepare error: %v\n", err)
		return nil, fmt.Errorf("Prepare error: %v", err)
	}
	defer stmt.Close()

	var account uint64
	err = stmt.QueryRowContext(ctx, userId).Scan(&account)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		fmt.Printf("UserPgRepo.GetUserBalance(). QueryRowContext error: %v\n", err)
		return nil, err
	}

	return &model.DBUser{Id: userId, Account: account}, nil
}

// Пополняет баланс пользователя в БД
func (repo *UserPgRepo) TopUpUserBalance(userId uint64, amount uint64) (*model.DBUser, error) {

	// незачем идти в БД ради изменения суммы баланса на дельта ноль
	if amount == 0 {
		fmt.Println("UserPgRepo.TopUpUserBalance(). Args error: amount = 0")
		return nil, fmt.Errorf("Args error: amount = 0")
	}

	user, err := repo.checkIsExistUser(userId)
	if err != nil {
		fmt.Printf("UserPgRepo.TopUpUserBalance(). checkIsExistUser(): %v", err)
		return nil, fmt.Errorf("checkIsExistUser(): %v", err)
	}

	isNewUser := false
	// создадим пользователя, если еще нет
	if user == nil {
		_, err = repo.createUser(userId, "", "", amount)
		if err != nil {
			fmt.Printf("UserPgRepo.TopUpUserBalance(). createUser(): %v", err)
			return nil, err
		}
		isNewUser = true
	}

	ctx := context.Background()

	tx, err := repo.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		fmt.Printf("UserPgRepo.TopUpUserBalance(). BeginTx error: %v", err)
		return nil, fmt.Errorf("BeginTx error: %v", err)
	}

	// процесс записи итогов транзакции, в таблицу транзакций
	stmt, err := tx.Prepare("UPDATE users SET balance = balance + $1 WHERE user_id = $2")
	if err != nil {
		fmt.Printf("UserPgRepo.TopUpUserBalance(). Prepare error: %v", err)
		_ = tx.Rollback()
		return nil, fmt.Errorf("Prepare error: %v", err)
	}
	defer stmt.Close()

	// Обновим баланс пользователя в БД
	res, err := stmt.Exec(amount, userId) // TODO Проверить транзакцию

	// ошибка при запросе к БД
	if err != nil {
		fmt.Printf("UserPgRepo.TopUpUserBalance(). Exec error: %v", err)
		_ = tx.Rollback()
		return nil, fmt.Errorf("Exec error: %v", err)
	}

	// ошибка при чтении полученных строк
	row, err := res.RowsAffected()
	if err != nil {
		fmt.Printf("UserPgRepo.TopUpUserBalance(). RowsAffected error: %v\r\n", err)
		_ = tx.Rollback()
		return nil, fmt.Errorf("RowsAffected error: %v", err)
	}

	if row != 1 {
		fmt.Println("UserPgRepo.TopUpUserBalance(). RowsAffected != 1. User not found")
		_ = tx.Rollback()
		return nil, fmt.Errorf("User not found")
	}

	// процесс записи итогов транзакции, в таблицу транзакций
	stmt, err = tx.Prepare("INSERT INTO transactions (user_id, amount, operation, date) VALUES ($1, $2, $3, $4)")
	if err != nil {
		fmt.Printf("UserPgRepo.TopUpUserBalance(). Prepare error: %v\r\n", err)
		_ = tx.Rollback()
		return nil, fmt.Errorf("Prepare error: %v", err)
	}
	defer stmt.Close()

	date := time.Now().Format("01-02-2006 15:04:05") // mm-dd-yy hh-mm-ss
	operation := "Top-up by transfer " + " +" + strconv.FormatUint(amount, 10) + " rub"
	_, err = stmt.Exec(userId, amount, operation, date)
	if err != nil {
		fmt.Printf("UserPgRepo.TopUpUserBalance(). Exec error: %v\r\n", err)
		_ = tx.Rollback()
		return nil, fmt.Errorf("Exec error: %v", err)
	}

	if err := tx.Commit(); err != nil {
		fmt.Printf("UserPgRepo.TopUpUserBalance(). Commit error: %v\r\n", err)
		return nil, fmt.Errorf("Commit error: %v", err)
	}

	// Вернем нового пользователя с балансом, который только что пополнили
	if isNewUser {
		return &model.DBUser{Id: userId, Account: amount}, nil
	}

	// Если пользователь существовал, вернем его с новым балансом
	user, err = repo.GetUserBalance(userId)
	if err != nil {
		fmt.Printf("UserPgRepo.TopUpUserBalance(). GetUserBalance error: %v\r\n", err)
		return nil, fmt.Errorf("GetUserBalance error: %v", err)
	}

	return user, nil
}

// Списывает средства со счета пользователя в БД
func (repo *UserPgRepo) DebitUserBalance(userId uint64, amount uint64) (*model.DBUser, error) {

	// незачем идти в БД ради изменения суммы баланса на дельта ноль
	if amount == 0 {
		fmt.Println("UserPgRepo.DebitUserBalance(). Args error: amount = 0")
		return nil, fmt.Errorf("Args error: amount = 0")
	}

	user, err := repo.checkIsExistUser(userId)
	if err != nil {
		fmt.Printf("UserPgRepo.DebitUserBalance(). checkIsExistUser(): %v\r\n", err)
		return nil, fmt.Errorf("checkIsExistUser: %v", err)
	}

	if user == nil {
		fmt.Println("UserPgRepo.DebitUserBalance(). User not found")
		return nil, fmt.Errorf("User not found")
	}

	ctx := context.Background()

	if user.Account < amount {
		fmt.Println("UserPgRepo.DebitUserBalance(). User has no balance")
		return nil, fmt.Errorf("User has no balance")
	}

	tx, err := repo.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		fmt.Printf("UserPgRepo.DebitUserBalance(). BeginTx error: %v\r\n", err)
		return nil, fmt.Errorf("BeginTx error: %v", err)
	}

	// процесс записи итогов транзакции, в таблицу транзакций
	stmt, err := tx.Prepare("UPDATE users SET balance = balance - $1 WHERE user_id = $2")
	if err != nil {
		fmt.Printf("UserPgRepo.DebitUserBalance(). Prepare error: %v\r\n", err)
		_ = tx.Rollback()
		return nil, fmt.Errorf("Prepare error: %v", err)
	}
	defer stmt.Close()

	// Обновим баланс пользователя в БД
	res, err := stmt.Exec(amount, userId) // TODO Проверить транзакцию
	// ошибка при запросе к БД
	if err != nil {
		fmt.Printf("UserPgRepo.DebitUserBalance(). Exec error: %v\r\n", err)
		_ = tx.Rollback()
		return nil, fmt.Errorf("Exec error: %v", err)
	}

	// ошибка при чтении полученных строк
	row, err := res.RowsAffected()
	if err != nil {
		fmt.Printf("UserPgRepo.DebitUserBalance(). RowsAffected error: %v\r\n", err)
		_ = tx.Rollback()
		return nil, fmt.Errorf("RowsAffected error: %v", err)
	}

	if row != 1 {
		fmt.Println("UserPgRepo.DebitUserBalance(). RowsAffected != 1. User not found")
		_ = tx.Rollback()
		return nil, fmt.Errorf("User not found")
	}

	// процесс записи итогов транзакции, в таблицу транзакций
	stmt, err = tx.Prepare("INSERT INTO transactions (user_id, amount, operation, date) VALUES ($1, $2, $3, $4)")
	if err != nil {
		fmt.Printf("UserPgRepo.DebitUserBalance(). Prepare error: %v\r\n", err)
		_ = tx.Rollback()
		return nil, fmt.Errorf("Prepare error: %v", err)
	}
	defer stmt.Close()

	date := time.Now().Format("01-02-2006 15:04:05") // mm-dd-yy hh-mm-ss
	operation := "Debit by transfer " + " -" + strconv.FormatUint(amount, 10) + " rub"
	_, err = stmt.Exec(userId, amount, operation, date)
	if err != nil {
		fmt.Printf("UserPgRepo.DebitUserBalance(). Exec: %v\r\n", err)
		_ = tx.Rollback()
		return nil, fmt.Errorf("Exec error: %v", err)
	}

	if err := tx.Commit(); err != nil {
		fmt.Printf("UserPgRepo.DebitUserBalance(). Commit error: %v\r\n", err)
		return nil, fmt.Errorf("Commit error: %v", err)
	}

	return &model.DBUser{Id: userId, Account: (user.Account - amount)}, nil
}

// Пересылает средства от одного пользователя к другому
func (repo *UserPgRepo) Transfer(userFromKey, userToKey, amount uint64) (*model.DBUser, error) {

	if amount <= 0 {
		fmt.Println("UserPgRepo.Transfer(). Args error: amount = 0")
		return nil, fmt.Errorf("Args error: amount = 0")
	}

	ctx := context.Background()
	tx, err := repo.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		fmt.Printf("UserPgRepo.Transfer(). BeginTx error: %v\r\n", err)
		return nil, fmt.Errorf("BeginTx error: %v", err)
	}

	stmt, err := tx.Prepare("UPDATE users SET balance = balance - $1 WHERE user_id = $2")
	if err != nil {
		fmt.Printf("UserPgRepo.Transfer(). Prepare error: %v\r\n", err)
		return nil, fmt.Errorf("Prepare error: %v", err)
	}
	defer stmt.Close()

	// спишем деньги со счета отправителя
	res, err := stmt.Exec(amount, userFromKey) // TODO -amount Неработает преобразование в отрицательное число

	if err == sql.ErrNoRows {
		fmt.Println("UserPgRepo.Transfer(). User not found")
		_ = tx.Rollback()
		return nil, fmt.Errorf("User not found")
	}

	if err != nil {
		fmt.Printf("UserPgRepo.Transfer(). Exec(amount, userFromKey) error: %v\r\n", err)
		_ = tx.Rollback()
		return nil, fmt.Errorf("Exec error: %v", err)
	}

	resAff, errAff := res.RowsAffected()

	if errAff != nil {
		fmt.Printf("UserPgRepo.Transfer(). RowsAffected error: %v\r\n", err)
		_ = tx.Rollback()
		return nil, fmt.Errorf("RowAff error: %v", err)
	}

	if resAff != 1 {
		fmt.Println("UserPgRepo.Transfer(). RowsAffected != 1. User not found")
		_ = tx.Rollback()
		return nil, fmt.Errorf("User with user_id: %d not found", userFromKey)
	}

	stmt, err = tx.Prepare("UPDATE users SET balance = balance + $1 WHERE user_id = $2")
	if err != nil {
		fmt.Printf("UserPgRepo.Transfer(). Prepare error: %v\r\n", err)
		return nil, fmt.Errorf("Prepare error: %v", err)
	}
	defer stmt.Close()

	// зачислим деньги на счет получателя
	res, err = stmt.Exec(amount, userToKey)

	if err == sql.ErrNoRows {
		fmt.Println("UserPgRepo.Transfer(). User not found")
		_ = tx.Rollback()
		return nil, fmt.Errorf("User with user_id: %d not found", userToKey)
	}

	if err != nil {
		fmt.Printf("UserPgRepo.Transfer(). Exec(amount, userToKey) error: %v\r\n", err)
		_ = tx.Rollback()
		return nil, fmt.Errorf("DB error: %v", err)
	}

	resAff, errAff = res.RowsAffected()
	if errAff != nil {
		fmt.Printf("UserPgRepo.Transfer(). RowsAffected error: %v\r\n", err)
		_ = tx.Rollback()
		return nil, fmt.Errorf("RowAff error: %v", err)
	}

	if resAff != 1 {
		fmt.Println("UserPgRepo.Transfer(). RowsAffected != 1. User not found")
		_ = tx.Rollback()
		return nil, fmt.Errorf("User with user_id: %d not found", userToKey)
	}

	// процесс записи итогов транзакции, в таблицу транзакций
	stmt, err = tx.Prepare("INSERT INTO transactions (user_id, amount, operation, date) VALUES ($1, $2, $3, $4)")
	if err != nil {
		fmt.Printf("UserPgRepo.Transfer(). Prepare error: %v\r\n", err)
		_ = tx.Rollback()
		return nil, fmt.Errorf("Prepare error: %v", err)
	}
	defer stmt.Close()

	date := time.Now().Format("01-02-2006 15:04:05") // mm-dd-yy hh-mm-ss
	amountString := strconv.FormatUint(amount, 10)

	operation := "Debit by transfer " + " -" + amountString + " rub"
	_, err = stmt.Exec(userFromKey, amount, operation, date)
	if err != nil {
		fmt.Printf("UserPgRepo.Transfer(). Exec error: %v\r\n", err)
		_ = tx.Rollback()
		return nil, fmt.Errorf("Exec error: %v", err)
	}

	operation = "Top-up by transfer " + " +" + amountString + " rub"
	_, err = stmt.Exec(userToKey, amount, operation, date)
	if err != nil {
		fmt.Printf("UserPgRepo.Transfer(). Exec error: %v\r\n", err)
		_ = tx.Rollback()
		return nil, fmt.Errorf("Exec error: %v", err)
	}

	if err := tx.Commit(); err != nil {
		fmt.Printf("UserPgRepo.Transfer(). Commit error: %v\r\n", err)
		return nil, fmt.Errorf("Commit error")
	}

	userTo, err := repo.checkIsExistUser(userToKey)
	if err != nil {
		return nil, err
	}
	return userTo, nil
}

// Создает пользователя
func (repo *UserPgRepo) createUser(userId uint64, fName, lName string, amount uint64) (*model.DBUser, error) {
	ctx := context.Background()

	stmt, err := repo.db.PrepareContext(ctx, "INSERT INTO users (userId, firstName, lastName, balance) VALUES ($1, $2, $3, $4)")
	if err != nil {
		fmt.Printf("UserPgRepo.createUser(). Prepare error: %v\n", err)
		return nil, fmt.Errorf("Prepare error: %v", err)
	}
	defer stmt.Close()

	if err != nil {
		return nil, err
	}

	ans, err := stmt.Exec(userId, fName, lName, amount)
	if err != nil {
		return nil, fmt.Errorf("UserPgRepo.createUser(). Exec error: %v", err)
	}

	if r, _ := ans.RowsAffected(); r > 0 {
		fmt.Println("INSERT NEW USER OK") // TODO CHECK
	}

	return &model.DBUser{
		Id:      userId,
		Account: 0,
	}, nil
}

// Создает пользователя
func (repo *UserPgRepo) checkIsExistUser(userId uint64) (*model.DBUser, error) {
	ctx := context.Background()

	stmt, err := repo.db.PrepareContext(ctx, "SELECT * FROM users WHERE user_id = $1 FOR UPDATE")
	if err != nil {
		fmt.Printf("UserPgRepo.checkIsExistUser(). PrepareContext error: %v\n", err)
		return nil, fmt.Errorf("PrepareContext error: %v", err)
	}
	defer stmt.Close()

	user := model.DBUser{}
	err = stmt.QueryRowContext(ctx, userId).Scan(&user.Id, &user.FirstName, &user.LastName, &user.Account)
	if err == sql.ErrNoRows {
		fmt.Println("UserPgRepo.checkIsExistUser(). User not found")
		return nil, fmt.Errorf("User not found")
	}

	if err != nil {
		fmt.Printf("UserPgRepo.checkIsExistUser(). QueryRowContext error: %v\n", err)
		return nil, fmt.Errorf("QueryRowContext error: %v", err)
	}

	return &user, nil
}
