package services

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type MessageKind int

const (
	PlaceBid MessageKind = iota
)

type Message struct {
	Message string
	Kind    MessageKind
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
	defer func() {
		close(room.Broadcast)
		close(room.Register)
		close(room.Unregister)
	}()

	for {
		select {
		case client := <-room.Register:
			continue // r.RegisterClient(client)
		case client := <-room.Unregister:
			continue // room.UnregisterClient(client)
		case message := <-room.Broadcast:
			continue // room.BroadCastMessage(message)
		case <-room.Context.Done():
			fmt.Println("Auction ending")
			// notificar usuarios que o leilao acabou
		}
	}
}

func (room *AuctionRoom) RegisterClient(client *Client) {

}

func NewAuctionRoom(ctx context.Context, BidService BidServices, id uuid.UUID) *AuctionRoom {
	return &AuctionRoom{
		Id:          id,
		Broadcast:   make(chan Message),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Context:     ctx,
		BidServices: BidService,
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
