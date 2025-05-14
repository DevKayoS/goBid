package services

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type MessageKind int

const (
	// Request
	PlaceBid MessageKind = iota

	// Ok/sucess
	SuccessfullyPlaceBid

	// Errors
	FailedToPlaceBid
	InvalidJSON

	//info
	NewBidPlace
	AuctionFinished
)

type Message struct {
	Message string      `json:"message,omitempty"`
	Kind    MessageKind `json:"kind"`
	Amount  float64     `json:"amount,omitempty"`
	UserId  uuid.UUID   `json:"user_id,omitempty"`
}

type AuctionLobby struct {
	sync.Mutex
	Rooms map[uuid.UUID]*AuctionRoom
}

type AuctionRoom struct {
	Id          uuid.UUID
	Context     context.Context
	Broadcast   chan Message
	Register    chan *Client
	Unregister  chan *Client
	Clients     map[uuid.UUID]*Client
	BidServices *BidServices
}

func (room *AuctionRoom) Run() {
	slog.Info("Auction has begun", "auctionId", room.Id)
	defer func() {
		close(room.Broadcast)
		close(room.Register)
		close(room.Unregister)
	}()

	for {
		select {
		case client := <-room.Register:
			room.registerClient(client)
		case client := <-room.Unregister:
			room.unregisterClient(client)
		case message := <-room.Broadcast:
			room.broadCastMessage(message)
		case <-room.Context.Done():
			slog.Info("Auction has ended", "auctionId", room.Id)
			for _, client := range room.Clients {
				client.Send <- Message{Kind: AuctionFinished, Message: "ayctuin has been finished"}
			}
			return
		}
	}
}

func (room *AuctionRoom) registerClient(client *Client) {
	slog.Info("New User Connected", "Client", client)
	room.Clients[client.UserId] = client
}

func (room *AuctionRoom) unregisterClient(client *Client) {
	slog.Info("User disconnected", "Client", client)
	delete(room.Clients, client.UserId)
}

func (room *AuctionRoom) broadCastMessage(message Message) {
	slog.Info("New message recieved", "RoomId", room.Id, "message", message.Message, "user_Id", message.UserId)
	switch message.Kind {
	case PlaceBid:
		bid, err := room.BidServices.PlaceBid(room.Context, room.Id, message.UserId, message.Amount)
		if err != nil {
			if errors.Is(err, ErrBidIsTooLow) {
				if client, ok := room.Clients[message.UserId]; ok {
					client.Send <- Message{Kind: FailedToPlaceBid, Message: ErrBidIsTooLow.Error(), UserId: message.UserId}
				}
				return
			}
		}

		if Client, ok := room.Clients[message.UserId]; ok {
			Client.Send <- Message{Kind: SuccessfullyPlaceBid, Message: "Your bid was successfully placed", UserId: message.UserId}
		}

		for id, client := range room.Clients {
			newBidMessage := Message{Kind: NewBidPlace, Message: "A new bid was placed", Amount: bid.BidAmount, UserId: message.UserId}
			if id == message.UserId {
				continue
			}
			client.Send <- newBidMessage
		}
	case InvalidJSON:
		client, ok := room.Clients[message.UserId]
		if !ok {
			slog.Info("Client not found in a hasmap", "user_id", message.UserId)
			return
		}

		client.Send <- message
	}
}

func NewAuctionRoom(ctx context.Context, BidService BidServices, id uuid.UUID) *AuctionRoom {
	return &AuctionRoom{
		Id:          id,
		Broadcast:   make(chan Message),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Clients:     make(map[uuid.UUID]*Client),
		Context:     ctx,
		BidServices: &BidService,
	}
}

type Client struct {
	Room   *AuctionRoom
	Conn   *websocket.Conn
	UserId uuid.UUID
	Send   chan Message
}

func NewClient(room *AuctionRoom, conn *websocket.Conn, userId uuid.UUID) *Client {
	return &Client{
		Room:   room,
		Conn:   conn,
		UserId: userId,
		Send:   make(chan Message, 512), // pode receber ate 512 mensagens por vez
	}
}

const (
	MAX_MESSAGE_SIZE = 512
	READ_DEAD_LINE   = 60 * time.Second // 1 minuto
	PING_PERIOD      = (READ_DEAD_LINE * 9) / 10
	WRITE_WAIT       = 10 * time.Second
)

func (client *Client) ReadEventLoop() {
	defer func() {
		client.Room.Unregister <- client
		client.Conn.Close()
	}()

	client.Conn.SetReadLimit(MAX_MESSAGE_SIZE)
	client.Conn.SetReadDeadline(time.Now().Add(READ_DEAD_LINE))

	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(READ_DEAD_LINE))
		return nil
	})

	for {
		var message Message
		message.UserId = client.UserId

		err := client.Conn.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error("Unexpected close error", "error", err)
				return
			}

			client.Room.Broadcast <- Message{
				Kind:    InvalidJSON,
				Message: "this message should be a valid json",
				UserId:  message.UserId,
			}
			continue
		}
		client.Room.Broadcast <- message
	}
}

func (client *Client) WriteEventLoop() {
	ticker := time.NewTicker(PING_PERIOD)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			if !ok {
				client.Conn.WriteJSON(Message{
					Kind:    websocket.CloseMessage,
					Message: "Closing websocket connection",
				})
				return
			}

			if message.Kind == AuctionFinished {
				close(client.Send)
				return
			}

			err := client.Conn.WriteJSON(message)
			if err != nil {
				client.Room.Unregister <- client
			}

		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(WRITE_WAIT))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				slog.Error("Unexpected write error", "error", err)
				return
			}
		}
	}
}
