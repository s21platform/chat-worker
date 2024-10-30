package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/s21platform/chat-worker/internal/model"
)

type Server struct {
	upgrade   websocket.Upgrader
	clients   map[*model.Client]struct{}
	usecase   Usecase
	broadcast chan []byte
}

func New(usecase Usecase) *Server {
	upgrade := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	clients := make(map[*model.Client]struct{})
	broadcast := make(chan []byte)
	return &Server{
		upgrade:   upgrade,
		clients:   clients,
		usecase:   usecase,
		broadcast: broadcast,
	}
}

func (s *Server) ConnectClient(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Println("failed to upgrade websocket connection:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	var clientWS model.Message
	err = conn.ReadJSON(&clientWS)
	log.Println("connect:", clientWS.ClientID)
	if err != nil {
		log.Println("failed to read client:", err)
		return
	}

	client := &model.Client{
		Conn:     conn,
		Send:     make(chan []byte, 256),
		ClientId: clientWS.ClientID,
	}
	s.clients[client] = struct{}{}

	go func() {
		err := s.usecase.HandleMessages(r.Context(), client)
		if err != nil {
			log.Println("failed to handle messages:", err)
			delete(s.clients, client)
		}
	}()

	if err := s.usecase.WriteMessage(r.Context(), client); err != nil {
		log.Println("failed to write message:", err)
		delete(s.clients, client)
	}
}

func (s *Server) sendMessages() {
	for {
		msg := <-s.broadcast
		var message *model.Message
		err := json.Unmarshal(msg, &message)
		if err != nil {
			log.Println("failed to unmarshal message:", err)
			continue
		}
		for client := range s.clients {
			if client.ClientId != message.ClientID {
				select {
				case client.Send <- []byte(message.Content):
				default:
					delete(s.clients, client)
					close(client.Send)
				}
			}
		}
	}
}

func (s *Server) Run(port string) {
	http.HandleFunc("/", s.ConnectClient)

	go s.sendMessages()
	go s.usecase.ReadMessage(context.Background(), s.broadcast)

	log.Println("server listening")
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		panic(err)
	}
}
