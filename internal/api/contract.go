package api

import (
	"context"
	"github.com/s21platform/chat-worker/internal/model"
)

type Usecase interface {
	HandleMessages(ctx context.Context, client *model.Client) error
	WriteMessage(ctx context.Context, client *model.Client) error
	ReadMessage(ctx context.Context, broadcast chan []byte)
}
