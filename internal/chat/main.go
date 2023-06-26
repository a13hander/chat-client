package chat

import (
	"context"
	"time"

	chatV1 "github.com/a13hander/chat-server/pkg/chat_v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ Client = (*client)(nil)

type Message struct {
	From      string
	Text      string
	CreatedAt time.Time
}

type Client interface {
	CreateChat(ctx context.Context, usernames []string) (string, error)
	ConnectChat(ctx context.Context, chatID string, username string) (chatV1.ChatV1_ConnectChatClient, error)
	SendMessage(ctx context.Context, chatID string, message *Message) error
}

type client struct {
	client chatV1.ChatV1Client
}

func NewClient(cc *grpc.ClientConn) *client {
	return &client{
		client: chatV1.NewChatV1Client(cc),
	}
}

func (c *client) CreateChat(ctx context.Context, usernames []string) (string, error) {
	res, err := c.client.CreateChat(ctx, &chatV1.CreateChatRequest{
		Username: usernames,
	})
	if err != nil {
		return "", err
	}

	return res.GetChatId(), nil
}

func (c *client) ConnectChat(ctx context.Context, chatID string, username string) (chatV1.ChatV1_ConnectChatClient, error) {
	return c.client.ConnectChat(ctx, &chatV1.ConnectChatRequest{
		ChatId:   chatID,
		Username: username,
	})
}

func (c *client) SendMessage(ctx context.Context, chatID string, message *Message) error {
	_, err := c.client.SendMessage(ctx, &chatV1.SendMessageRequest{
		ChatId: chatID,
		Message: &chatV1.Message{
			From:      message.From,
			Text:      message.Text,
			CreatedAt: timestamppb.New(message.CreatedAt),
		},
	})
	if err != nil {
		return err
	}

	return nil
}
