package app

import (
	"context"
	"log"
	"time"

	authV1 "github.com/a13hander/auth-service-api/pkg/auth_v1"
	"github.com/a13hander/chat-client/internal/auth"
	"github.com/a13hander/chat-client/internal/chat"
	"github.com/a13hander/chat-client/internal/config"
	"github.com/a13hander/chat-client/internal/handler"
	"github.com/a13hander/chat-client/internal/interceptor"
	"github.com/a13hander/chat-client/internal/redis"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type ServiceProvider struct {
	authClient     auth.Client
	chatClient     chat.Client
	redisClient    redis.Client
	handlerService *handler.Handler
}

func NewServiceProvider() *ServiceProvider {
	return &ServiceProvider{}

}

func (s *ServiceProvider) GetAuthClient(ctx context.Context) auth.Client {
	if s.authClient == nil {
		creds, err := credentials.NewClientTLSFromFile("service.pem", "")
		if err != nil {
			log.Fatalf("could not process the credentials: %v", err)
		}

		conn, err := grpc.DialContext(ctx, config.GetConfig().AuthAddr, grpc.WithTransportCredentials(creds))
		if err != nil {
			log.Fatalf("could not process the credentials: %v", err)
		}

		closer.add(conn.Close)

		client := authV1.NewAuthV1Client(conn)
		s.authClient = auth.NewClient(client)
	}

	return s.authClient
}

func (s *ServiceProvider) GetChatClient(ctx context.Context) chat.Client {
	if s.chatClient == nil {
		authInterceptor := interceptor.NewAuthInterceptor(s.GetAuthClient(ctx), s.GetRedisClient())
		authInterceptor.Run(60*time.Minute, 1*time.Minute)

		conn, err := grpc.DialContext(
			ctx,
			config.GetConfig().ChatSeverAddr,
			grpc.WithUnaryInterceptor(authInterceptor.Unary),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			log.Fatalf("failed to dial chat service: %s", err.Error())
		}
		closer.add(func() error {
			if conn != nil {
				return conn.Close()
			}
			return nil
		})

		s.chatClient = chat.NewClient(conn)
	}

	return s.chatClient
}

func (s *ServiceProvider) GetRedisClient() redis.Client {
	if s.redisClient == nil {
		client := redis.NewClient(config.GetConfig().RedisAddr)
		closer.add(func() error {
			if client != nil {
				return client.Close()
			}
			return nil
		})

		err := client.Ping()
		if err != nil {
			log.Fatalf("failed to ping redis: %s", err.Error())
		}

		s.redisClient = client
	}

	return s.redisClient
}

func (s *ServiceProvider) GetHandlerService(ctx context.Context) *handler.Handler {
	if s.handlerService == nil {
		s.handlerService = handler.NewHandler(
			s.GetRedisClient(),
			s.GetAuthClient(ctx),
			s.GetChatClient(ctx),
		)
	}

	return s.handlerService
}
