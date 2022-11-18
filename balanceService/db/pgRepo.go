package db

import (
	"balance_service/pkg/mErrors"
	"balance_service/pkg/struct4parse"
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"log"
	"time"
)

const serviceIncome = 0
const orderIncome = 0

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type PostgresDB struct {
	db *sql.DB
}

func NewPostgresDB(cfg Config) (*PostgresDB, error) {
	url := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=%v",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.SSLMode)
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("postgres connect error : (%v)", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &PostgresDB{db: db}, nil
}

//AddIncome пополнение счета в бд
func (db *PostgresDB) AddIncome(ctx context.Context, income struct4parse.BalanceWithDesc) error {
	tx, err := db.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx failed: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		`INSERT INTO balance.history
				(user_id, service_id, order_id, value, occurred_at, description, replenish)
			VALUES
				($1, $2, $3, $4, $5, $6, $7)`,
		income.UserId, serviceIncome, orderIncome, income.Value, income.Time, "income "+income.Description, true)

	if err != nil {
		return fmt.Errorf("add transaction to history query exec failed: %w", err)
	}

	var isUserIdExist bool

	err = tx.QueryRow(
		"SELECT EXISTS(SELECT user_id FROM balance.balance WHERE user_id = $1) AS exists",
		income.UserId).Scan(&isUserIdExist)

	if err != nil {
		return fmt.Errorf("check user_id exists query row failed: %w", err)
	}
	if isUserIdExist { //если пользователь существует, то добавляем ему деньги
		_, err = tx.Exec(
			"UPDATE balance.balance SET value = value + $1 WHERE user_id = $2",
			income.Value, income.UserId)
		if err != nil {
			return fmt.Errorf("add income query exec failed: %w", err)
		}
	} else { //иначе создаем пользователя с данной в запросе суммой
		_, err = tx.Exec(
			"INSERT INTO balance.balance (user_id, value) VALUES($1, $2)",
			income.UserId, income.Value)
		if err != nil {
			return fmt.Errorf("add new user_id with balance query exec failed: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("tx commit failed failed: %w", err)
	}
	return nil
}

//AddReserve резервирует деньги для дальнейшего списания
func (db *PostgresDB) AddReserve(ctx context.Context, expense struct4parse.Transaction) error {
	tx, err := db.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx failed: %w", err)
	}
	defer tx.Rollback()
	var isUserIdExist bool

	err = tx.QueryRow(
		"SELECT EXISTS(SELECT user_id FROM balance.balance WHERE user_id = $1 ) AS exists",
		expense.UserId).Scan(&isUserIdExist)

	if err != nil {
		return fmt.Errorf("check user_id exists query row failed: %w", err)
	}
	if isUserIdExist {
		_, err = tx.Exec( //обновляем баланс: вычитаем зарезервированные деньги, если хватает
			"UPDATE balance.balance SET value = value - $1 WHERE user_id = $2",
			expense.Value, expense.UserId)
		if err != nil {
			return fmt.Errorf("user_id %d: %w", expense.UserId, mErrors.NotEnoughUserBalanceError)
		}
	} else {
		return fmt.Errorf("user_id %d: %w", expense.UserId, mErrors.UnknownUserIdError)
	}

	//добавляем деньги в таблицу резервации
	_, err = db.db.Query(
		`INSERT INTO balance.reserved
				(user_id, service_id, order_id, value )
			VALUES
				($1, $2, $3, $4) `,
		expense.UserId, expense.ServiceId, expense.OrderId, expense.Value)

	if err != nil {
		return fmt.Errorf("add transaction to history query exec failed: %w", err)
	}
	time := time.Now()
	_, err = tx.Exec(
		`INSERT INTO balance.history
				(user_id, service_id, order_id, value, occurred_at, description, replenish)
			VALUES
				($1, $2, $3, $4, $5, $6, $7)`,
		expense.UserId, expense.ServiceId, expense.OrderId, expense.Value, time, "reserve", false)

	if err != nil {
		return fmt.Errorf("add transaction to history query exec failed: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("tx commit failed failed: %w", err)
	}
	return nil
}

//AddExpense переводит деньги из резервации если они были, если нет то с основного счета пользователя.
//Выдает ошибку если денег нехватает.
//После успешной операции добавляет ее в отчет - транзакции
func (db *PostgresDB) AddExpense(ctx context.Context, expense struct4parse.Transaction) error {
	tx, err := db.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx failed: %w", err)
	}
	defer tx.Rollback()

	var isUserIdExist bool
	var index int64
	//проверка на существование резерва денег
	err = tx.QueryRow(
		"SELECT EXISTS(SELECT user_id FROM balance.reserved WHERE user_id = $1 and service_id = $2 and order_id = $3 and value = $4) AS exists",
		expense.UserId, expense.ServiceId, expense.OrderId, expense.Value).Scan(&isUserIdExist)

	if err != nil {
		return fmt.Errorf("check user_id exists query row failed: %w", err)
	}
	if isUserIdExist {
		row := tx.QueryRow("SELECT id FROM balance.reserved WHERE user_id = $1 and service_id = $2 and order_id = $3 and value = $4 ",
			expense.UserId, expense.ServiceId, expense.OrderId, expense.Value)
		err = row.Scan(&index)
		_, err = db.db.Exec(
			"DELETE FROM balance.reserved WHERE id = $1", index) //если было в резерве - удаляем из него
		if err != nil {

			if errPq, ok := err.(*pgconn.PgError); ok {
				if errPq.Code == pgerrcode.CheckViolation {
					return fmt.Errorf("user_id %d: %w", expense.UserId, mErrors.NotEnoughUserBalanceError)
				}
			}

			return fmt.Errorf("add expense query exec failed: %w", err)
		}
	} else {
		//пытаемся списать с основного счета если в резерве не было
		err = tx.QueryRow(
			"SELECT EXISTS(SELECT user_id FROM balance.balance WHERE user_id = $1) AS exists",
			expense.UserId).Scan(&isUserIdExist)

		if err != nil {
			return fmt.Errorf("check user_id exists query row failed: %w", err)
		}
		if isUserIdExist {
			_, err = db.db.Exec(
				"UPDATE balance.balance SET value = value - $1 WHERE user_id = $2",
				expense.Value, expense.UserId)
			if err != nil {

				if errPq, ok := err.(*pgconn.PgError); ok {
					if errPq.Code == pgerrcode.CheckViolation {
						return fmt.Errorf("user_id %d: %w", expense.UserId, mErrors.NotEnoughUserBalanceError)
					}
				}

				return fmt.Errorf("add expense query exec failed: %w", err)
			}
		} else {
			return fmt.Errorf("user_id %d: %w", expense.UserId, mErrors.UnknownUserIdError)
		}

		if err = tx.Commit(); err != nil {
			return fmt.Errorf("tx commit failed failed: %w", err)
		}

	}
	//добавляем в отчет транзакцию
	_, err = db.db.Exec(
		`INSERT INTO balance.history
				(user_id, service_id, order_id, value, occurred_at, description, replenish)
			VALUES
				($1, $2, $3, $4, $5, $6, $7)`,
		expense.UserId, expense.ServiceId, expense.OrderId, expense.Value, expense.Time, expense.Description, false)
	if err != nil {
		return fmt.Errorf("add transaction to history query exec failed: %w", err)
	}
	return nil
}

//GetBalance узнаем баланс одного пользователя
func (db *PostgresDB) GetBalance(ctx context.Context, income *struct4parse.Balance) error {
	tx, err := db.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx failed: %w", err)
	}
	defer tx.Rollback()

	var isUserIdExist bool
	var balanceValue float64

	err = tx.QueryRow(
		"SELECT EXISTS(SELECT user_id FROM balance.balance WHERE user_id = $1) AS exists",
		income.UserId).Scan(&isUserIdExist)

	if err != nil {
		return fmt.Errorf("check user_id exists query row failed: %w", err)
	}

	if isUserIdExist {
		err = tx.QueryRow(
			"SELECT value FROM balance.balance WHERE user_id = $1", income.UserId).Scan(&balanceValue)
		if err != nil {
			return fmt.Errorf("get balance query row failed: %w", err)
		}
	} else {
		return fmt.Errorf("user_id %d: %w", income.UserId, mErrors.UnknownUserIdError)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("tx commit failed failed: %w", err)
	}
	income.Value = balanceValue
	return nil
}

//GetAllBalances получаем список балансов всех пользователей
func (db *PostgresDB) GetAllBalances(ctx context.Context, income *[]struct4parse.Balance) error {
	tx, err := db.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx failed: %w", err)
	}
	defer tx.Rollback()

	rows, err := db.db.Query("select user_id, value from balance.balance")
	defer rows.Close()

	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("get balance query row failed: %s", err.Error())
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("tx commit failed failed: %s", err.Error())
	}

	var elem struct4parse.Balance
	arr := make([]struct4parse.Balance, 0)
	for rows.Next() {
		err := rows.Scan(&elem.UserId, &elem.Value)
		if err != nil {
			return fmt.Errorf("err in red rows: %s", err)
		}
		arr = append(arr, elem)
	}
	*income = arr
	err = rows.Err()
	return nil
}

//GetAllReserved получаем список всех резервов
func (db *PostgresDB) GetAllReserved(ctx context.Context, income *[]struct4parse.Reserve) error {
	tx, err := db.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx failed: %w", err)
	}
	defer tx.Rollback()

	rows, err := db.db.Query("select * from balance.reserved")
	if err != nil {
		return fmt.Errorf("get balance query row failed: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("tx commit failed failed: %w", err)
	}
	defer rows.Close()

	var elem struct4parse.Reserve
	arr := make([]struct4parse.Reserve, 0)
	for rows.Next() {
		err := rows.Scan(&elem.Id, &elem.UserId, &elem.ServiceId, &elem.OrderId, &elem.Value)
		if err != nil {
			fmt.Errorf("err in red rows in Reserved: %s", err.Error())
		}
		arr = append(arr, elem)
	}

	*income = arr
	err = rows.Err()
	return nil
}

//GetAllTransactions получаем список балансов всех транзакций
func (db *PostgresDB) GetAllTransactions(ctx context.Context, income *[]struct4parse.Transaction) error {
	tx, err := db.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx failed: %w", err)
	}
	defer tx.Rollback()

	rows, err := db.db.Query("select * from balance.history")
	defer rows.Close()

	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("get balance query row failed: %s", err.Error())
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("tx commit failed failed: %s", err.Error())
	}

	var elem struct4parse.Transaction
	arr := make([]struct4parse.Transaction, 0)
	for rows.Next() {
		err := rows.Scan(&elem.Id, &elem.UserId, &elem.ServiceId, &elem.OrderId, &elem.Value, &elem.Time, &elem.Description, &elem.Replenish)
		if err != nil {
			return fmt.Errorf("err in red rows: %s", err)
		}
		arr = append(arr, elem)
	}
	*income = arr
	err = rows.Err()
	return nil
}

func (db *PostgresDB) DisReserve(ctx context.Context, expense struct4parse.Transaction) error {
	tx, err := db.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx failed: %w", err)
	}
	defer tx.Rollback()
	var isUserIdExist bool

	err = tx.QueryRow(
		"SELECT EXISTS(SELECT user_id FROM balance.reserved WHERE user_id = $1 and service_id = $2 and order_id = $3 and value = $4) AS exists",
		expense.UserId, expense.ServiceId, expense.OrderId, expense.Value).Scan(&isUserIdExist)

	if err != nil {
		return fmt.Errorf("check user_id exists query row failed: %w", err)
	}
	if isUserIdExist {
		_, err = tx.Exec( //обновляем баланс: вычитаем зарезервированные деньги, если хватает
			"UPDATE balance.balance SET value = value + $1 WHERE user_id = $2",
			expense.Value, expense.UserId)
		if err != nil {
			return fmt.Errorf("add reserve query exec failed: %w", err)
		}
	} else {
		return fmt.Errorf("service_id: %d, order_id:%d err: %w", expense.ServiceId, expense.OrderId, mErrors.UnknownReserveError)
	}

	var index int64
	if expense.Id == 0 {
		row := tx.QueryRow("SELECT id FROM balance.reserved WHERE user_id = $1 and service_id = $2 and order_id = $3 and value = $4 ",
			expense.UserId, expense.ServiceId, expense.OrderId, expense.Value)
		err = row.Scan(&index)
	} else {
		index = expense.Id
	}

	//удаляем деньги из таблицы резервации
	_, err = db.db.Exec(
		"DELETE FROM balance.reserved WHERE id = $1", index)
	time := time.Now()
	if err != nil {

		if errPq, ok := err.(*pgconn.PgError); ok {
			if errPq.Code == pgerrcode.CheckViolation {
				return fmt.Errorf("user_id %d: %w", expense.UserId, mErrors.NotEnoughUserBalanceError)
			}
		}

		return fmt.Errorf("add expense query exec failed: %w", err)
	}
	_, err = tx.Exec(
		`INSERT INTO balance.history
				(user_id, service_id, order_id, value, occurred_at, description, replenish)
			VALUES
				($1, $2, $3, $4, $5, $6, $7)`,
		expense.UserId, expense.ServiceId, expense.OrderId, expense.Value, time, "disreserve", true)

	if err != nil {
		return fmt.Errorf("add transaction to history query exec failed: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("tx commit failed failed: %w", err)
	}
	return nil
}
