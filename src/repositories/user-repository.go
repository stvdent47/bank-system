package repositories

import (
	"bank-system/src/entities"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) (string, error)
	FindByUsername(ctx context.Context, username string) (entities.User, error)
}

type UserRepositoryPgx struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepositoryPgx {
	return &UserRepositoryPgx{pool: pool}
}

func (this *UserRepositoryPgx) Create(ctx context.Context, user *entities.User) (string, error) {
	_, err := this.pool.Exec(
		ctx,
		"insert into users (id, username, password, email) values ($1, $2, $3, $4)",
		user.ID,
		user.Username,
		user.Password,
		user.Email,
	)

	if err != nil {
		return "", err
	}

	return user.ID, nil
}

func (this *UserRepositoryPgx) FindByUsername(ctx context.Context, username string) (entities.User, error) {
	user := entities.User{}

	this.pool.QueryRow(
		ctx,
		"select id, username, password, email from users where username = $1",
		username,
	).Scan(&user.ID, &user.Username, &user.Password, &user.Email)

	return user, nil
}
