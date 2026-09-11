package postgres

import (
	"context"
	"database/sql"

	"be_evindo/internal/entity/user"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repository *UserRepository) FindByID(ctx context.Context, id int64) (*user.User, error) {
	result := &user.User{}
	err := repository.db.QueryRowContext(ctx,
		`SELECT id, name, email, password FROM users WHERE id = $1`, id,
	).Scan(&result.ID, &result.Name, &result.Email, &result.Password)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (repository *UserRepository) Create(ctx context.Context, entity *user.User) error {
	row := repository.db.QueryRowContext(ctx,
		`INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id`,
		entity.Name, entity.Email, entity.Password,
	)
	return row.Scan(&entity.ID)
}
