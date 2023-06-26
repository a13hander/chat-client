package handler

import (
	"github.com/a13hander/chat-client/internal/auth"
	"github.com/a13hander/chat-client/internal/chat"
	"github.com/a13hander/chat-client/internal/redis"
)

type Handler struct {
	redisClient redis.Client
	authClient  auth.Client
	chatClient  chat.Client
}

func NewHandler(redisClient redis.Client, authClient auth.Client, chatClient chat.Client) *Handler {
	return &Handler{
		redisClient: redisClient,
		authClient:  authClient,
		chatClient:  chatClient,
	}
}
