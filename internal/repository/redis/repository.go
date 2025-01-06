package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/s21platform/chat-worker/internal/config"
)

type Repo struct {
	conn    *redis.Client
	channel string
}

func MustConnect(cfg *config.Config) *Repo {
	conn := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
	})
	return &Repo{
		conn:    conn,
		channel: cfg.Redis.Channel,
	}
}

func (r *Repo) Close() {
	_ = r.conn.Close()
}

func (r *Repo) PublishMessage(ctx context.Context, msg []byte) error {
	err := r.conn.Publish(ctx, r.channel, msg)
	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}
	return nil
}

func (r *Repo) HandleMessagesFromChannel(ctx context.Context) *redis.PubSub {
	return r.conn.Subscribe(ctx, r.channel)
}
