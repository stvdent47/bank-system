package repositories

import (
	"bank-system/src/entities"
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentRepository interface {
	Create(ctx context.Context, data entities.Payment) (string, error)
	CreateMany(ctx context.Context, data []entities.Payment) ([]string, error)
	GetByLoanId(ctx context.Context, id string) ([]entities.Payment, error)
	Update(ctx context.Context, data entities.Payment) (string, error)
	GetOverdue(ctx context.Context) ([]entities.Payment, error)
}

type PaymentRepositoryPgx struct {
	pool *pgxpool.Pool
}

func NewPaymentRepository(pool *pgxpool.Pool) *PaymentRepositoryPgx {
	return &PaymentRepositoryPgx{pool: pool}
}

func (this *PaymentRepositoryPgx) Create(ctx context.Context, data entities.Payment) (string, error) {
	_, err := this.pool.Exec(
		ctx,
		"insert into payments (id, loan_id, amount, due_date, principal_part, interest_part, is_paid) values ($1, $2, $3, $4, $5, $6, $7, $8)",
		data.ID,
		data.LoanId,
		data.Amount,
		data.DueDate,
		data.PrincipalPart,
		data.InterestPart,
		data.IsPaid,
	)

	if err != nil {
		return "", err
	}

	return "", nil
}

func (this *PaymentRepositoryPgx) CreateMany(ctx context.Context, data []entities.Payment) ([]string, error) {
	tx, err := this.pool.Begin(ctx)
	if err != nil {
		return []string{}, err
	}
	defer tx.Rollback(ctx)

	batch := new(pgx.Batch)

	for _, payment := range data {
		batch.Queue(
			"insert into payments (id, loan_id, amount, due_date, principal_part, interest_part, is_paid) values ($1, $2, $3, $4, $5, $6, $7)",
			payment.ID,
			payment.LoanId,
			payment.Amount,
			payment.DueDate,
			payment.PrincipalPart,
			payment.InterestPart,
			payment.IsPaid,
		)
	}

	result := tx.SendBatch(ctx, batch)

	err = result.Close()
	if err != nil {
		return []string{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return []string{}, err
	}

	paymentIds := make([]string, len(data))
	for i, payment := range data {
		paymentIds[i] = payment.ID
	}

	return paymentIds, nil
}

func (this *PaymentRepositoryPgx) GetByLoanId(ctx context.Context, id string) ([]entities.Payment, error) {
	rows, err := this.pool.Query(
		ctx,
		"select id, loan_id, amount, date, due_date, principal_part, interest_part, is_paid from payments where loan_id = $1",
		id,
	)
	if err != nil {
		return []entities.Payment{}, err
	}
	defer rows.Close()

	var payments []entities.Payment
	for rows.Next() {
		var payment entities.Payment

		var nullableDate *int64

		err = rows.Scan(
			&payment.ID,
			&payment.LoanId,
			&payment.Amount,
			&nullableDate,
			&payment.DueDate,
			&payment.PrincipalPart,
			&payment.InterestPart,
			&payment.IsPaid,
		)
		if nullableDate != nil {
			payment.Date = *nullableDate
		}

		if err != nil {
			return []entities.Payment{}, err
		}

		payments = append(payments, payment)
	}

	return payments, nil
}

func (this *PaymentRepositoryPgx) Update(ctx context.Context, data entities.Payment) (string, error) {
	_, err := this.pool.Exec(
		ctx,
		"update payments set amount = $1, date = $2, due_date = $3, principal_part = $4, interest_part = $5, status = $6, is_paid = $7 where id = $1",
		data.Amount,
		data.Date,
		data.DueDate,
		data.PrincipalPart,
		data.InterestPart,
		data.Status,
		data.IsPaid,
	)

	if err != nil {
		return "", err
	}

	return data.ID, nil
}

func (this *PaymentRepositoryPgx) GetOverdue(ctx context.Context) ([]entities.Payment, error) {
	// Get current timestamp in Unix format
	currentTime := time.Now().Unix()

	rows, err := this.pool.Query(
		ctx,
		`select id, loan_id, amount, date, due_date, principal_part, interest_part, status, is_paid 
			from payments 
		where due_date < $1 and is_paid = false and (status = 'new')`,
		currentTime,
	)
	if err != nil {
		return []entities.Payment{}, err
	}
	defer rows.Close()

	var payments []entities.Payment
	for rows.Next() {
		var payment entities.Payment
		var nullableDate *int64

		err = rows.Scan(
			&payment.ID,
			&payment.LoanId,
			&payment.Amount,
			&nullableDate,
			&payment.DueDate,
			&payment.PrincipalPart,
			&payment.InterestPart,
			&payment.Status,
			&payment.IsPaid,
		)

		if err != nil {
			return []entities.Payment{}, err
		}

		if nullableDate != nil {
			payment.Date = *nullableDate
		}

		payments = append(payments, payment)
	}

	return payments, nil
}
