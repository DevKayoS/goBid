package services

import (
	"context"
	"errors"

	"github.com/DevKayoS/goBid/internal/store/pgstore"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicatedEmailOrUsername = errors.New("username or email already exists")
	ErrInvalidCredentials        = errors.New("Invalid Credentials")
)

type UserServices struct {
	pool    *pgxpool.Pool
	queries *pgstore.Queries
}

func NewUserService(pool *pgxpool.Pool) UserServices {
	return UserServices{
		pool:    pool,
		queries: pgstore.New(pool),
	}
}

func (us *UserServices) CreateUser(ctx context.Context, userName, email, password, bio string) (uuid.UUID, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return uuid.UUID{}, err
	}

	args := pgstore.CreateUserParams{
		UserName:     userName,
		Email:        email,
		PasswordHash: hash,
		Bio:          bio,
	}

	id, err := us.queries.CreateUser(ctx, args)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgstore.UNIQUE_VALUE {
			return uuid.UUID{}, ErrDuplicatedEmailOrUsername
		}

		return uuid.UUID{}, err
	}

	return id, nil
}

func (us *UserServices) AuthenticateUser(ctx context.Context, email, password string) (uuid.UUID, error) {
	user, err := us.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.UUID{}, ErrInvalidCredentials
		}

		return uuid.UUID{}, err
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return uuid.UUID{}, ErrInvalidCredentials
		}
		return uuid.UUID{}, err
	}

	return user.ID, nil
}
