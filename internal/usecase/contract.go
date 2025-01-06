package usecase

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type Redis interface {
	PublishMessage(ctx context.Context, msg []byte) error
	HandleMessagesFromChannel(ctx context.Context) *redis.PubSub
}
