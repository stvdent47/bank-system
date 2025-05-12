package db

import (
	"bank-system/config"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

func createUserTable(db *pgxpool.Pool, ctx context.Context) error {
	_, err := db.Exec(
		ctx,
		`create table if not exists users (
			id varchar(100) primary key,
			username varchar(50) not null unique,
			password varchar(255) not null,
			email varchar(50) not null unique
		)`,
	)

	return err
}

func createTransactionTable(db *pgxpool.Pool, ctx context.Context) error {
	_, err := db.Exec(
		ctx,
		`create table if not exists transactions (
			id varchar(100) primary key,
			amount bigint not null,
			from_id varchar(100),
			to_id varchar(100),
			type varchar(50) not null check (type in ('deposit', 'payment', 'transfer', 'withdrawal')),
			description varchar(255) not null,
			created_at bigint not null
		)`,
	)

	return err
}

func createAccountTable(db *pgxpool.Pool, ctx context.Context) error {
	_, err := db.Exec(
		ctx,
		`create table if not exists accounts (
			id varchar(100) primary key,
			balance bigint not null,
			user_id varchar(100) not null references users(id) on delete cascade,
			created_at bigint not null
		)`,
	)

	return err
}

func createCardTable(db *pgxpool.Pool, ctx context.Context) error {
	_, err := db.Exec(
		ctx,
		`create table if not exists cards (
			id varchar(100) primary key,
			number bytea not null unique,
			expiration bytea not null,
			cvv varchar(255) not null,
			user_id varchar(100) not null references users(id) on delete cascade,
			account_id varchar(100) not null references accounts(id) on delete cascade,
			created_at bigint not null
		)`,
	)

	return err
}

func createPaymentTable(db *pgxpool.Pool, ctx context.Context) error {
	_, err := db.Exec(
		ctx,
		`create table if not exists payments (
			id varchar(100) primary key,
			loan_id varchar(100) not null references loans(id) on delete cascade,
			amount bigint not null,
			date bigint,
			due_date bigint not null,
			principal_part bigint not null,
			interest_part bigint not null,
			status varchar(50) not null check (status in ('new', 'paid', 'overdue')),
			is_paid boolean not null
		)`,
	)

	return err
}

func createLoanTable(db *pgxpool.Pool, ctx context.Context) error {
	_, err := db.Exec(
		ctx,
		`create table if not exists loans (
			id varchar(100) primary key,
			user_id varchar(100) not null references users(id) on delete cascade,
			account_id varchar(100) not null references accounts(id) on delete cascade,
			amount bigint not null,
			interest_rate double precision not null,
			term int not null,
			start_date bigint not null,
			debt bigint not null
		)`,
	)

	return err
}

func createTables(db *pgxpool.Pool, ctx context.Context) error {
	err := createUserTable(db, ctx)
	if err != nil {
		return err
	}

	err = createTransactionTable(db, ctx)
	if err != nil {
		return err
	}

	err = createAccountTable(db, ctx)
	if err != nil {
		return err
	}

	err = createCardTable(db, ctx)
	if err != nil {
		return err
	}

	err = createPaymentTable(db, ctx)
	if err != nil {
		return err
	}

	err = createLoanTable(db, ctx)
	if err != nil {
		return err
	}

	return nil
}

func New(ctx context.Context, logger *logrus.Logger) (*pgxpool.Pool, error) {
	config := config.LoadDBConfig()

	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
	)

	db, err := pgxpool.New(ctx, url)

	if err != nil {
		return nil, err
	}

	err = db.Ping(ctx)
	if err != nil {
		return nil, err
	}

	err = createTables(db, ctx)
	if err != nil {
		return nil, err
	}

	logger.Info("Successfully connected to db")

	return db, nil
}
