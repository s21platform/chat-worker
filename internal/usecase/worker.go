package usecase

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/s21platform/chat-worker/internal/model"
	"log"
)

type Usecase struct {
	r Redis
}

func New(r Redis) *Usecase {
	return &Usecase{r: r}
}

func (u *Usecase) HandleMessages(ctx context.Context, client *model.Client) error {
	for {
		_, msg, err := client.Conn.ReadMessage()
		if err != nil {
			return fmt.Errorf("failed to read message: %w", err)
		}

		if err := u.r.PublishMessage(ctx, msg); err != nil {
			log.Printf("failed to publish message: %v", err)
		}
	}
}

func (u *Usecase) WriteMessage(ctx context.Context, client *model.Client) error {
	for msg := range client.Send {
		if err := client.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			return fmt.Errorf("failed to write message: %w", err)
		}
	}
	return nil
}

func (u *Usecase) ReadMessage(ctx context.Context, broadcast chan []byte) {
	pubsub := u.r.HandleMessagesFromChannel(ctx)
	defer pubsub.Close()

	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			log.Printf("failed to receive message: %v", err)
			return
		}

		broadcast <- []byte(msg.Payload)
	}
}
