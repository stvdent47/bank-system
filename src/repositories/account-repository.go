package repositories

import (
	"bank-system/src/entities"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountRepository interface {
	Create(ctx context.Context, data entities.Account) (string, error)
	GetAll(ctx context.Context, userId string) ([]entities.Account, error)
	GetById(ctx context.Context, id string) (entities.Account, error)
	UpdateBalance(ctx context.Context, id string, newBalance int64) error
	Delete(ctx context.Context, id string) (string, error)
	Transfer(ctx context.Context, senderId string, receiverId string, amount int64) error
}

type AccountRepositoryPgx struct {
	pool *pgxpool.Pool
}

func NewAccountRepository(pool *pgxpool.Pool) *AccountRepositoryPgx {
	return &AccountRepositoryPgx{pool: pool}
}

func (this *AccountRepositoryPgx) Create(ctx context.Context, data entities.Account) (string, error) {
	_, err := this.pool.Exec(
		ctx,
		"insert into accounts (id, balance, user_id, created_at) values ($1, $2, $3, $4)",
		data.ID,
		data.Balance,
		data.UserID,
		data.CreatedAt,
	)

	if err != nil {
		return "", err
	}

	return data.ID, nil
}

func (this *AccountRepositoryPgx) GetAll(ctx context.Context, userId string) ([]entities.Account, error) {
	rows, err := this.pool.Query(
		ctx,
		"select id, balance, user_id, created_at from accounts where user_id = $1",
		userId,
	)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	var accounts []entities.Account

	for rows.Next() {
		var account entities.Account
		err := rows.Scan(
			&account.ID,
			&account.Balance,
			&account.UserID,
			&account.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	if len(accounts) == 0 {
		return nil, nil
	}

	return accounts, nil
}

func (this *AccountRepositoryPgx) GetById(ctx context.Context, id string) (entities.Account, error) {
	row := this.pool.QueryRow(
		ctx,
		"select id, balance, user_id, created_at from accounts where id = $1",
		id,
	)

	var account entities.Account
	err := row.Scan(
		&account.ID,
		&account.Balance,
		&account.UserID,
		&account.CreatedAt,
	)

	if err != nil {
		return entities.Account{}, err
	}

	return account, nil
}

func (this *AccountRepositoryPgx) UpdateBalance(ctx context.Context, id string, newBalance int64) error {
	_, err := this.pool.Exec(
		ctx,
		"update accounts set balance = $1 where id = $2",
		newBalance,
		id,
	)

	if err != nil {
		return err
	}

	return nil
}

func (this *AccountRepositoryPgx) Delete(ctx context.Context, id string) (string, error) {
	_, err := this.pool.Exec(
		ctx,
		"delete from accounts where id = $1",
		id,
	)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (this *AccountRepositoryPgx) Transfer(ctx context.Context, fromId string, toId string, amount int64) error {
	tx, err := this.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(
		ctx,
		"update accounts set balance = balance - $1 where id = $2",
		amount,
		fromId,
	)
	if err != nil {
		println("first")
		return err
	}

	_, err = tx.Exec(
		ctx,
		"update accounts set balance = balance + $1 where id = $2",
		amount,
		toId,
	)
	if err != nil {
		println("second")
		return err
	}

	tx.Commit(ctx)

	return nil
}
