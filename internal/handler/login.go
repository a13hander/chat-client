package handler

import (
	"context"

	"github.com/a13hander/chat-client/internal/redis"
)

func (h *Handler) Login(ctx context.Context, username, password string) error {
	refreshToken, err := h.authClient.GetRefreshToken(ctx, username, password)
	if err != nil {
		return err
	}

	accessToken, err := h.authClient.GetAccessToken(ctx, refreshToken)
	if err != nil {
		return err
	}

	err = h.redisClient.Set(redis.UsernameKey, username, 0)
	if err != nil {
		return err
	}

	err = h.redisClient.Set(redis.PasswordKey, password, 0)
	if err != nil {
		return err
	}

	err = h.redisClient.Set(redis.AccessTokenKey, accessToken, 0)
	if err != nil {
		return err
	}

	err = h.redisClient.Set(redis.RefreshTokenKey, refreshToken, 0)
	if err != nil {
		return err
	}

	return nil
}
