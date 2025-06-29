package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/leodanuarta/go-grpc-ecommerce-be/internal/entity"
)

func (ar *authRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	row := ar.db.QueryRowContext(ctx, "SELECT id, email, password, full_name, role_code FROM \"user\" WHERE email = $1 AND is_deleted IS false", email)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var user entity.User
	err := row.Scan(
		&user.Id,
		&user.Email,
		&user.Password,
		&user.FullName,
		&user.RoleCode,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, err
}

func (ar *authRepository) InsertUser(ctx context.Context, user *entity.User) error {
	_, err := ar.db.ExecContext(
		ctx,
		"INSERT INTO \"user\" (id, full_name, email, password, role_code, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, is_deleted) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
		user.Id,
		user.FullName,
		user.Email,
		user.Password,
		user.RoleCode,
		user.CreatedAt,
		user.CreatedBy,
		user.UpdatedAt,
		user.UpdatedBy,
		user.DeletedAt,
		user.DeletedBy,
		user.IsDeleted,
	)

	if err != nil {
		return err
	}

	return nil
}
