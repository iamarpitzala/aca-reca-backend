package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/config"
	"github.com/redis/go-redis/v9"
)

type SessionService struct {
	rdb *redis.Client
	ttl time.Duration
}

type SessionData struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	SessionID uuid.UUID `json:"session_id"`
	UserAgent string    `json:"user_agent"`
	IPAddress string    `json:"ip_address"`
	CreatedAt time.Time `json:"created_at"`
}

func NewSessionService(rdb *redis.Client, cfg config.SessionConfig) *SessionService {
	return &SessionService{
		rdb: rdb,
		ttl: cfg.TTL,
	}
}

func (ss *SessionService) StoreSession(ctx context.Context, sessionID uuid.UUID, data *SessionData) error {
	key := fmt.Sprintf("session:%s", sessionID.String())

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal session data: %w", err)
	}

	if err := ss.rdb.Set(ctx, key, jsonData, ss.ttl).Err(); err != nil {
		return fmt.Errorf("failed to store session in Redis: %w", err)
	}

	return nil
}

func (ss *SessionService) GetSession(ctx context.Context, sessionID uuid.UUID) (*SessionData, error) {
	key := fmt.Sprintf("session:%s", sessionID.String())

	val, err := ss.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("session not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get session from Redis: %w", err)
	}

	var data SessionData
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session data: %w", err)
	}

	return &data, nil
}

func (ss *SessionService) DeleteSession(ctx context.Context, sessionID uuid.UUID) error {
	key := fmt.Sprintf("session:%s", sessionID.String())

	if err := ss.rdb.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete session from Redis: %w", err)
	}

	return nil
}

func (ss *SessionService) RefreshSession(ctx context.Context, sessionID uuid.UUID) error {
	key := fmt.Sprintf("session:%s", sessionID.String())

	if err := ss.rdb.Expire(ctx, key, ss.ttl).Err(); err != nil {
		return fmt.Errorf("failed to refresh session TTL: %w", err)
	}

	return nil
}

func (ss *SessionService) StoreUserSessions(ctx context.Context, userID uuid.UUID, sessionIDs []uuid.UUID) error {
	key := fmt.Sprintf("user:sessions:%s", userID.String())

	sessionStrs := make([]string, len(sessionIDs))
	for i, id := range sessionIDs {
		sessionStrs[i] = id.String()
	}

	jsonData, err := json.Marshal(sessionStrs)
	if err != nil {
		return fmt.Errorf("failed to marshal session IDs: %w", err)
	}

	if err := ss.rdb.Set(ctx, key, jsonData, ss.ttl).Err(); err != nil {
		return fmt.Errorf("failed to store user sessions in Redis: %w", err)
	}

	return nil
}

func (ss *SessionService) GetUserSessions(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	key := fmt.Sprintf("user:sessions:%s", userID.String())

	val, err := ss.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return []uuid.UUID{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user sessions from Redis: %w", err)
	}

	var sessionStrs []string
	if err := json.Unmarshal([]byte(val), &sessionStrs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session IDs: %w", err)
	}

	sessionIDs := make([]uuid.UUID, len(sessionStrs))
	for i, str := range sessionStrs {
		id, err := uuid.Parse(str)
		if err != nil {
			continue
		}
		sessionIDs[i] = id
	}

	return sessionIDs, nil
}
