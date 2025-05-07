package api

import (
	"github.com/DevKayoS/goBid/internal/services"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

type Api struct {
	Router         *chi.Mux
	UserService    services.UserServices
	Sessions       *scs.SessionManager
	ProductService services.ProductService
	BidService     services.BidServices
	WsUpgrader     websocket.Upgrader
	AuctionLobby   services.AuctionLobby
}
