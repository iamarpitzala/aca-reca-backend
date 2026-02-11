package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) port.UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO tbl_user (id, email, password, first_name, last_name, phone, created_at, updated_at)
		VALUES (:id, :email, :password, :first_name, :last_name, :phone, :created_at, :updated_at)`
	_, err := r.db.NamedExecContext(ctx, query, user)
	return err
}

func (r *userRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `SELECT id, email, password, first_name, last_name, phone, created_at, updated_at FROM tbl_user WHERE id = $1 AND deleted_at IS NULL`
	var user domain.User
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("failed to find user")
	}
	return &user, nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, email, password, first_name, last_name, phone, created_at, updated_at FROM tbl_user WHERE LOWER(email) = LOWER($1) AND deleted_at IS NULL`
	var user domain.User
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.New("failed to get user by email")
	}
	return &user, nil
}

func (r *userRepo) EmailExists(ctx context.Context, email string) (bool, error) {
	query := `SELECT COUNT(*) FROM tbl_user WHERE LOWER(email) = LOWER($1) AND deleted_at IS NULL`
	var count int
	err := r.db.GetContext(ctx, &count, query, email)
	if err != nil {
		return false, errors.New("failed to check if email exists")
	}
	return count > 0, nil
}

func (r *userRepo) Update(ctx context.Context, user *domain.User) error {
	query := `UPDATE tbl_user SET first_name = :first_name, last_name = :last_name, phone = :phone, updated_at = :updated_at WHERE id = :id`
	_, err := r.db.NamedExecContext(ctx, query, user)
	return err
}

type sessionRepo struct {
	db *sqlx.DB
}

func NewSessionRepository(db *sqlx.DB) port.SessionRepository {
	return &sessionRepo{db: db}
}

func (r *sessionRepo) Create(ctx context.Context, session *domain.Session) error {
	query := `INSERT INTO tbl_session (id, user_id, refresh_token, expires_at, created_at, updated_at)
		VALUES (:id, :user_id, :refresh_token, :expires_at, :created_at, :updated_at)`
	_, err := r.db.NamedExecContext(ctx, query, session)
	return err
}

func (r *sessionRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Session, error) {
	query := `SELECT id, user_id, refresh_token, expires_at, created_at, updated_at FROM tbl_session WHERE id = $1 AND deleted_at IS NULL`
	var session domain.Session
	err := r.db.GetContext(ctx, &session, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("session not found")
		}
		return nil, errors.New("failed to find session")
	}
	return &session, nil
}

func (r *sessionRepo) GetByRefreshToken(ctx context.Context, refreshToken string) (*domain.Session, error) {
	query := `SELECT id, user_id, refresh_token, expires_at, created_at, updated_at FROM tbl_session WHERE refresh_token = $1 AND deleted_at IS NULL`
	var session domain.Session
	err := r.db.GetContext(ctx, &session, query, refreshToken)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("session not found")
		}
		return nil, errors.New("failed to find session")
	}
	return &session, nil
}

func (r *sessionRepo) Update(ctx context.Context, session *domain.Session) error {
	query := `UPDATE tbl_session SET refresh_token = :refresh_token, expires_at = :expires_at, updated_at = :updated_at WHERE id = :id`
	_, err := r.db.NamedExecContext(ctx, query, session)
	return err
}

func (r *sessionRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE tbl_session SET deleted_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

func (r *sessionRepo) GetUserSessions(ctx context.Context, userID uuid.UUID) ([]domain.Session, error) {
	query := `SELECT id, user_id, refresh_token, expires_at, created_at, updated_at FROM tbl_session WHERE user_id = $1 AND deleted_at IS NULL`
	var sessions []domain.Session
	err := r.db.SelectContext(ctx, &sessions, query, userID)
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *sessionRepo) Revoke(ctx context.Context, sessionID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `UPDATE tbl_session SET updated_at = $1 WHERE id = $2`, time.Now(), sessionID)
	if err != nil {
		return errors.New("failed to revoke session")
	}
	return nil
}
