package main

import (
	"github.com/s21platform/chat-worker/internal/api"
	"github.com/s21platform/chat-worker/internal/config"
	"github.com/s21platform/chat-worker/internal/repository/redis"
	"github.com/s21platform/chat-worker/internal/usecase"
)

func main() {
	cfg := config.MustLoad()
	r := redis.MustConnect(cfg)

	uc := usecase.New(r)
	server := api.New(uc)
	server.Run()
}
