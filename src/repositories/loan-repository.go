package repositories

import (
	"bank-system/src/entities"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LoanRepository interface {
	Create(ctx context.Context, data entities.Loan) (string, error)
	GetById(ctx context.Context, id string) (entities.Loan, error)
}

type LoanRepositoryPgx struct {
	pool *pgxpool.Pool
}

func NewLoanRepository(pool *pgxpool.Pool) *LoanRepositoryPgx {
	return &LoanRepositoryPgx{pool: pool}
}

func (this *LoanRepositoryPgx) Create(ctx context.Context, data entities.Loan) (string, error) {
	_, err := this.pool.Exec(
		ctx,
		"insert into loans (id, user_id, account_id, amount, interest_rate, term, start_date, debt) values ($1, $2, $3, $4, $5, $6, $7, $8)",
		data.ID,
		data.UserId,
		data.AccountId,
		data.Amount,
		data.InterestRate,
		data.Term,
		data.StartDate,
		data.Debt,
	)

	if err != nil {
		return "", err
	}

	return data.ID, nil
}

func (this *LoanRepositoryPgx) GetById(ctx context.Context, id string) (entities.Loan, error) {
	row := this.pool.QueryRow(
		ctx,
		"select id, user_id, account_id, amount, interest_rate, term, start_date, debt from loans where id = $1",
		id,
	)

	var loan entities.Loan
	err := row.Scan(
		&loan.ID,
		&loan.UserId,
		&loan.AccountId,
		&loan.Amount,
		&loan.InterestRate,
		&loan.Term,
		&loan.StartDate,
		&loan.Debt,
	)

	if err != nil {
		return entities.Loan{}, err
	}

	return loan, nil
}
