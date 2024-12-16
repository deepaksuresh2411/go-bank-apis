package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(int, *CreateAccountReq) (*Account, error)
	GetAccountByID(int) (*Account, error)
	GetAccounts() ([]*Account, error)
}

type PostgresStorage struct {
	db *sql.DB
}

func (p *PostgresStorage) Init() error {
	return p.CreateAccountTable()
}

func NewPostgresStorage() (*PostgresStorage, error) {
	configStr := "user=deepak dbname=gobank password=openit sslmode=disable"
	db, err := sql.Open("postgres", configStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStorage{
		db: db,
	}, nil
}

func (p *PostgresStorage) CreateAccountTable() error {
	AccountTablequery := `create table if not exists account(
		id serial primary key,
		firstname varchar(25),
		lastname varchar(25),
		account_number serial,
		balance int,
		created_at timestamp
	)`
	_, err := p.db.Exec(AccountTablequery)
	return err
}

func (p *PostgresStorage) CreateAccount(a *Account) error {
	query := `insert into account 
	(firstname, lastname, account_number, balance, created_at)
	values ($1, $2, $3, $4, $5)`
	resp, err := p.db.Query(query,
		a.FirstName,
		a.LastName,
		a.AccountNumber,
		a.Balance,
		a.CreatedAt,
	)
	if err != nil {
		return err
	}
	fmt.Println(resp)
	return nil
}

func (p *PostgresStorage) DeleteAccount(id int) error {
	query := "delete from account where id=$1"
	_, err := p.db.Query(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresStorage) UpdateAccount(id int, acc *CreateAccountReq) (*Account, error) {
	account := &Account{}
	query := `update account set firstname=$2, lastname=$3 where id=$1`
	_, err := p.db.Query(query, id, acc.FirstName, acc.LastName)
	if err != nil {
		return nil, err
	}
	rows, err := p.db.Query("select * from account where id=$1", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err := rows.Scan(
			&account.ID,
			&account.FirstName,
			&account.LastName,
			&account.AccountNumber,
			&account.Balance,
			&account.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
	}
	return account, nil
}

func (p *PostgresStorage) GetAccountByID(id int) (*Account, error) {
	query := "select * from account where id=$1"
	rows, _ := p.db.Query(query, id)
	for rows.Next() {
		var acc = Account{}
		err := rows.Scan(
			&acc.ID,
			&acc.FirstName,
			&acc.LastName,
			&acc.AccountNumber,
			&acc.Balance,
			&acc.CreatedAt,
		)
		if err != nil {
			return nil, err
		} else {
			return &acc, nil
		}
	}
	return nil, fmt.Errorf("account with id %d not found", id)
}

func (p *PostgresStorage) GetAccounts() ([]*Account, error) {
	accounts := []*Account{}
	query := "select * from account"
	rows, _ := p.db.Query(query)
	for rows.Next() {
		account := Account{}
		err := rows.Scan(
			&account.ID,
			&account.FirstName,
			&account.LastName,
			&account.AccountNumber,
			&account.Balance,
			&account.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, &account)
	}
	return accounts, nil
}
