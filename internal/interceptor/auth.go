package interceptor

import (
	"context"
	"fmt"
	"time"

	"github.com/a13hander/chat-client/internal/auth"
	"github.com/a13hander/chat-client/internal/redis"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AuthInterceptor struct {
	authClient  auth.Client
	redisClient redis.Client
}

func NewAuthInterceptor(authClient auth.Client, redisClient redis.Client) *AuthInterceptor {
	return &AuthInterceptor{
		authClient:  authClient,
		redisClient: redisClient,
	}
}

func (i *AuthInterceptor) Run(refreshTokenPeriod time.Duration, accessTokenPeriod time.Duration) {
	go func() {
		t := time.NewTicker(refreshTokenPeriod)
		ctx := context.Background()

		for _ = range t.C {
			username, err := i.redisClient.Get(redis.UsernameKey)
			if err != nil {
				fmt.Println("failed to get username from redis")
				continue
			}

			password, err := i.redisClient.Get(redis.PasswordKey)
			if err != nil {
				fmt.Println("failed to get password from redis")
				continue
			}

			refreshToken, err := i.authClient.GetRefreshToken(ctx, username, password)
			if err != nil {
				fmt.Println("failed to get refresh token")
				continue
			}

			err = i.redisClient.Set(redis.RefreshTokenKey, refreshToken, 0)
			if err != nil {
				fmt.Println("failed to set refresh token to redis")
				continue
			}

			fmt.Println("refresh token has been updated")
		}
	}()

	go func() {
		t := time.NewTicker(accessTokenPeriod)
		ctx := context.Background()

		for _ = range t.C {
			refreshToken, err := i.redisClient.Get(redis.RefreshTokenKey)
			if err != nil {
				fmt.Println("failed to get refresh token from redis")
				continue
			}

			accessToken, err := i.authClient.GetAccessToken(ctx, refreshToken)
			if err != nil {
				fmt.Println("failed to get access token")
				continue
			}

			err = i.redisClient.Set(redis.AccessTokenKey, accessToken, 0)
			if err != nil {
				fmt.Println("failed to set access token to redis")
				continue
			}

			fmt.Println("access token has been updated")
		}
	}()
}

func (i *AuthInterceptor) Unary(ctx context.Context, method string, req interface{}, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	accessToken, err := i.redisClient.Get(redis.AccessTokenKey)
	if err != nil {
		return err
	}

	md := metadata.New(map[string]string{"Authorization": "Bearer " + accessToken})
	ctx = metadata.NewOutgoingContext(ctx, md)

	return invoker(ctx, method, req, reply, cc, opts...)
}
