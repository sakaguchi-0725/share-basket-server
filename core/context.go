package core

import (
	"context"
	"errors"
)

type contextKey string

const (
	UserIDKey contextKey = "userID"
	TxKey     contextKey = "transaction"
)

// GetUserID はコンテキストからユーザーIDを取得する。
// ユーザーIDが存在しない場合はUnauthorizedエラーを返す。
func GetUserID(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok || userID == "" {
		return "", NewAppError(ErrUnauthorized, errors.New("failed to retrieve authentication credentials"))
	}
	return userID, nil
}
