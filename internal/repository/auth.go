package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

func CreateUser(ctx context.Context, db *sqlx.DB, user *domain.User) error {
	query := `INSERT INTO tbl_user (id, email, password, first_name, last_name, phone, created_at, updated_at)
		VALUES (:id, :email, :password, :first_name, :last_name, :phone, :created_at, :updated_at)`
	_, err := db.NamedExecContext(ctx, query, user)
	if err != nil {
		return err
	}
	return nil
}

func GetUserByEmail(ctx context.Context, db *sqlx.DB, email string) (*domain.User, error) {
	query := `SELECT id, email, password, first_name, last_name, phone, created_at, updated_at FROM tbl_user WHERE email = $1 AND deleted_at IS NULL`
	var user domain.User
	err := db.GetContext(ctx, &user, query, email)
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.New("failed to get user by email")
	}
	return &user, nil
}

func EmailExists(ctx context.Context, db *sqlx.DB, email string) (bool, error) {
	query := `SELECT COUNT(*) FROM tbl_user WHERE email = $1 AND deleted_at IS NULL`
	var count int
	err := db.GetContext(ctx, &count, query, email)
	if err != nil {
		return false, errors.New("failed to check if email exists")
	}
	return count > 0, nil
}

func CreateSession(ctx context.Context, db *sqlx.DB, session *domain.Session) error {
	query := `INSERT INTO tbl_session (id, user_id, refresh_token, is_active, expires_at, created_at, updated_at)
		VALUES (:id, :user_id, :refresh_token, :is_active, :expires_at, :created_at, :updated_at)`

	_, err := db.NamedExecContext(ctx, query, session)
	if err != nil {
		return err
	}
	return nil
}

func GetSessionByRefreshToken(ctx context.Context, db *sqlx.DB, refreshToken string) (*domain.Session, error) {
	query := `SELECT id, user_id, refresh_token, is_active, expires_at, created_at, updated_at FROM tbl_session WHERE refresh_token = $1 AND deleted_at IS NULL`
	var session domain.Session
	err := db.GetContext(ctx, &session, query, refreshToken)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("session not found")
		}
		return nil, errors.New("failed to find session")
	}
	return &session, nil
}

func UpdateSession(ctx context.Context, db *sqlx.DB, session *domain.Session) error {
	query := `UPDATE tbl_session SET refresh_token = :refresh_token, is_active = :is_active, expires_at = :expires_at, updated_at = :updated_at WHERE id = :id`
	_, err := db.NamedExecContext(ctx, query, session)
	if err != nil {
		return err
	}
	return nil
}

func DeleteSession(ctx context.Context, db *sqlx.DB, sessionID uuid.UUID) error {
	query := `UPDATE tbl_session SET deleted_at = $1 WHERE id = $2`
	_, err := db.ExecContext(ctx, query, time.Now(), sessionID)
	if err != nil {
		return err
	}
	return nil
}

func GetSessionByID(ctx context.Context, db *sqlx.DB, sessionID uuid.UUID) (*domain.Session, error) {
	query := `SELECT id, user_id, refresh_token, is_active, expires_at, created_at, updated_at FROM tbl_session WHERE id = $1 AND deleted_at IS NULL`
	var session domain.Session
	err := db.GetContext(ctx, &session, query, sessionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("session not found")
		}
		return nil, errors.New("failed to find session")
	}
	return &session, nil
}

func GetUserByID(ctx context.Context, db *sqlx.DB, userID uuid.UUID) (*domain.User, error) {
	query := `SELECT id, email, password, first_name, last_name, phone, created_at, updated_at FROM tbl_user WHERE id = $1 AND deleted_at IS NULL`
	var user domain.User
	err := db.GetContext(ctx, &user, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("failed to find user")
	}
	return &user, nil
}

func UpdateUser(ctx context.Context, db *sqlx.DB, user *domain.User) error {
	query := `UPDATE tbl_user SET first_name = :first_name, last_name = :last_name, phone = :phone, updated_at = :updated_at WHERE id = :id`
	_, err := db.NamedExecContext(ctx, query, user)
	if err != nil {
		return err
	}
	return nil
}

func GetUserSessions(ctx context.Context, db *sqlx.DB, userId uuid.UUID) ([]domain.Session, error) {
	var sessions []domain.Session
	err := db.Select(&sessions,
		"SELECT id, user_id, refresh_token, is_active, expires_at, created_at, updated_at FROM tbl_session WHERE user_id = $1 AND is_active = $2 AND deleted_at IS NULL",
		userId, true)
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

func RevokeSession(ctx context.Context, db *sqlx.DB, sessionID uuid.UUID) error {
	_, err := db.Exec("UPDATE tbl_session SET is_active = $1, updated_at = $2 WHERE id = $3", false, time.Now(), sessionID)
	if err != nil {
		return errors.New("failed to revoke session")
	}

	return nil
}
