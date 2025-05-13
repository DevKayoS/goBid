package services

import (
	"context"
	"errors"
	"log/slog"
	"sync"

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

	//info
	NewBidPlace
	AuctionFinished
)

type Message struct {
	Message string
	Kind    MessageKind
	UserId  uuid.UUID
	Amount  float64
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
					client.Send <- Message{Kind: FailedToPlaceBid, Message: ErrBidIsTooLow.Error()}
				}
				return
			}
		}

		if Client, ok := room.Clients[message.UserId]; ok {
			Client.Send <- Message{Kind: SuccessfullyPlaceBid, Message: "Your bid was successfully placed"}
		}

		for id, client := range room.Clients {
			newBidMessage := Message{Kind: NewBidPlace, Message: "A new bid was placed", Amount: bid.BidAmount}
			if id == message.UserId {
				continue
			}
			client.Send <- newBidMessage
		}
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
