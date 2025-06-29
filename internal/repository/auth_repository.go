package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/leodanuarta/go-grpc-ecommerce-be/internal/entity"
)

func (ar *authRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	row := ar.db.QueryRowContext(ctx, "SELECT id, email, password, full_name, role_code, created_at FROM \"user\" WHERE email = $1 AND is_deleted IS false", email)
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
		&user.CreatedAt,
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

func (ar *authRepository) UpdateUserPassword(ctx context.Context, userid string, hashedNewPassword string, updatedBy string) error {
	_, err := ar.db.ExecContext(
		ctx,
		"UPDATE \"user\" SET password = $1, updated_at = $2, updated_by = $3 WHERE id = $4",
		hashedNewPassword,
		time.Now(),
		updatedBy,
		userid,
	)

	if err != nil {
		return err
	}

	return nil
}
