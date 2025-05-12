package repositories

import (
	"bank-system/src/entities"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CardRepository interface {
	Create(ctx context.Context, data entities.Card) (string, error)
	GetById(ctx context.Context, id string) (entities.Card, error)
}

type CardRepositoryPgx struct {
	pool *pgxpool.Pool
}

func NewCardRepository(pool *pgxpool.Pool) *CardRepositoryPgx {
	return &CardRepositoryPgx{pool: pool}
}

func (this *CardRepositoryPgx) Create(ctx context.Context, data entities.Card) (string, error) {
	_, err := this.pool.Exec(
		ctx,
		"insert into cards (id, number, expiration, cvv, user_id, account_id, created_at) values ($1, $2, $3, $4, $5, $6, $7)",
		data.ID,
		data.Number,
		data.Expiration,
		data.CVV,
		data.UserId,
		data.AccountID,
		data.CreatedAt,
	)

	if err != nil {
		return "", err
	}

	return data.ID, nil
}

func (this *CardRepositoryPgx) GetById(ctx context.Context, id string) (entities.Card, error) {
	row := this.pool.QueryRow(
		ctx,
		"select id, number, expiration, cvv, user_id, account_id, created_at from cards where id = $1",
		id,
	)

	var card entities.Card
	err := row.Scan(
		&card.ID,
		&card.Number,
		&card.Expiration,
		&card.CVV,
		&card.UserId,
		&card.AccountID,
		&card.CreatedAt,
	)

	if err != nil {
		return entities.Card{}, err
	}

	return card, nil
}
