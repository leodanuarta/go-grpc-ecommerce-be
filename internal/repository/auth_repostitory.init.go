package repository

import (
	"context"
	"database/sql"

	"github.com/leodanuarta/go-grpc-ecommerce-be/internal/entity"
)

type IAuthRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	InsertUser(ctx context.Context, user *entity.User) error
	UpdateUserPassword(ctx context.Context, userid string, hashedNewPassword string, updatedBy string) error
}
type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) IAuthRepository {
	return &authRepository{
		db: db,
	}
}
