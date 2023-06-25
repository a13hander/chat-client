package auth

import (
	"context"

	authV1 "github.com/a13hander/auth-service-api/pkg/auth_v1"
)

var _ Client = (*client)(nil)

type Client interface {
	GetRefreshToken(ctx context.Context, username, password string) (string, error)
	GetAccessToken(ctx context.Context, refreshToken string) (string, error)
}

type client struct {
	client authV1.AuthV1Client
}

func NewClient(cl authV1.AuthV1Client) *client {
	return &client{
		client: cl,
	}
}

func (c *client) GetRefreshToken(ctx context.Context, username, password string) (string, error) {
	token, err := c.client.GetRefreshToken(ctx, &authV1.GetRefreshTokenRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return "", err
	}

	return token.GetRefreshToken(), nil
}

func (c *client) GetAccessToken(ctx context.Context, refreshToken string) (string, error) {
	token, err := c.client.GetAccessToken(ctx, &authV1.GetAccessTokenRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		return "", err
	}

	return token.GetAccessToken(), nil
}
