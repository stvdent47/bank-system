package repositories

import (
	"bank-system/src/entities"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository interface {
	Create(ctx context.Context, data entities.Transaction) (string, error)
	GetByAccountId(ctx context.Context, accountId string) ([]entities.Transaction, error)
}
type TransactionRepositoryPgx struct {
	pool *pgxpool.Pool
}

func NewTransactionRepository(pool *pgxpool.Pool) *TransactionRepositoryPgx {
	return &TransactionRepositoryPgx{pool: pool}
}

func (this *TransactionRepositoryPgx) Create(ctx context.Context, data entities.Transaction) (string, error) {
	_, err := this.pool.Exec(
		ctx,
		"insert into transactions (id, amount, from_id, to_id, type, description, created_at) values ($1, $2, $3, $4, $5, $6, $7)",
		data.ID,
		data.Amount,
		data.FromAccountId,
		data.ToAccountId,
		data.Type,
		data.Description,
		data.CreatedAt,
	)

	if err != nil {
		return "", err
	}

	return data.ID, nil
}

func (this *TransactionRepositoryPgx) GetByAccountId(ctx context.Context, id string) ([]entities.Transaction, error) {
	rows, err := this.pool.Query(
		ctx,
		"select id, amount, from_id, to_id, type, description, created_at from transactions where from_id = $1 or to_id = $1",
		id,
	)
	if err != nil {
		return []entities.Transaction{}, err
	}
	defer rows.Close()

	var transactions []entities.Transaction
	for rows.Next() {
		var transaction entities.Transaction

		var nullableTo *string

		err := rows.Scan(
			&transaction.ID,
			&transaction.Amount,
			&transaction.FromAccountId,
			// &transaction.ToAccountId,
			&nullableTo,
			&transaction.Type,
			&transaction.Description,
			&transaction.CreatedAt,
		)
		if nullableTo != nil {
			transaction.ToAccountId = *nullableTo
		}

		if err != nil {
			return []entities.Transaction{}, err
		}

		transactions = append(transactions, transaction)
	}
	return transactions, nil
}
