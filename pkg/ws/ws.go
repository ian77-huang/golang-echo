package ws

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v5"
)

type WebSocketPacket struct {
	ID   string
	Conn *websocket.Conn
}

type webSocketHub struct {
	clients    map[string]*websocket.Conn
	broadcast  chan []byte
	single     chan *MessagePacket
	register   chan *WebSocketPacket
	unregister chan *WebSocketPacket
	mu         sync.Mutex
}

type WebSocketHub interface {
	Run()
	New(c *echo.Context, ID string, callbackMessage func(messageType int, p []byte)) error
	Single(p *MessagePacket)
}

type MessagePacket struct {
	ID      string
	Message []byte
}

func upgrader() *websocket.Upgrader {
	return &websocket.Upgrader{}
}

func NewHub() WebSocketHub {
	return &webSocketHub{
		clients:    make(map[string]*websocket.Conn),
		broadcast:  make(chan []byte),
		single:     make(chan *MessagePacket),
		register:   make(chan *WebSocketPacket),
		unregister: make(chan *WebSocketPacket),
	}
}

func (h *webSocketHub) New(c *echo.Context, ID string, callbackMessage func(messageType int, p []byte)) error {
	conn, err := upgrader().Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	h.register <- &WebSocketPacket{ID: ID, Conn: conn}

	defer func() {
		h.unregister <- &WebSocketPacket{ID: ID, Conn: conn}
		conn.Close()
	}()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		callbackMessage(messageType, message)
	}

	return nil
}

// 監聽並處理連線註冊、註銷與廣播
func (h *webSocketHub) Run() {
	go func() {
		for {
			select {
			case conn := <-h.register:
				h.mu.Lock()
				h.clients[conn.ID] = conn.Conn
				h.mu.Unlock()
				log.Println("新用戶加入，目前連線數：", len(h.clients))

			case conn := <-h.unregister:
				h.mu.Lock()
				if _, ok := h.clients[conn.ID]; ok {
					delete(h.clients, conn.ID)
					conn.Conn.Close()
				}
				h.mu.Unlock()
				log.Println("用戶離開，目前連線數：", len(h.clients))

			case packet := <-h.single:
				h.mu.Lock()

				if conn, ok := h.clients[packet.ID]; ok {
					err := conn.WriteMessage(websocket.TextMessage, packet.Message)
					if err != nil {
						log.Println("發送失敗，關閉連線：", err)
					}
				}

				h.mu.Unlock()
			case message := <-h.broadcast:
				h.mu.Lock()
				// 將訊息發送給所有在線的 client
				for key, conn := range h.clients {

					err := conn.WriteMessage(websocket.TextMessage, message)
					if err != nil {
						log.Println("發送失敗，關閉連線：", err)
						conn.Close()
						delete(h.clients, key)
					}
				}
				h.mu.Unlock()
			}
		}
	}()
}

func (h *webSocketHub) Single(p *MessagePacket) {
	h.single <- p
}

// func (h *webSocketHub) Run() {
// 	for {
// 		select {
// 		case conn := <-h.register:
// 			h.mu.Lock()
// 			h.clients[conn.ID] = conn.Conn
// 			h.mu.Unlock()
// 			log.Println("新用戶加入，目前連線數：", len(h.clients))

// 		case conn := <-h.unregister:
// 			h.mu.Lock()
// 			if _, ok := h.clients[conn.ID]; ok {
// 				delete(h.clients, conn.ID)
// 				conn.Conn.Close()
// 			}
// 			h.mu.Unlock()
// 			log.Println("用戶離開，目前連線數：", len(h.clients))

// 		case message := <-h.broadcast:
// 			h.mu.Lock()
// 			// 將訊息發送給所有在線的 client
// 			for key, conn := range h.clients {

// 				err := conn.WriteMessage(websocket.TextMessage, message)
// 				if err != nil {
// 					log.Println("發送失敗，關閉連線：", err)
// 					conn.Close()
// 					delete(h.clients, key)
// 				}
// 			}
// 			h.mu.Unlock()
// 		}
// 	}
// }
