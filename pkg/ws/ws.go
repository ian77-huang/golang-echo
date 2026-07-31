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
	register   chan *WebSocketPacket
	unregister chan *WebSocketPacket
	mu         sync.Mutex
}

type WebSocketHub interface {
	Run()
	New(c *echo.Context) error
}

func NewWebSocketUpgrader() *websocket.Upgrader {
	return &websocket.Upgrader{}
}

func NewWebSocketHub() WebSocketHub {
	return &webSocketHub{
		clients:    make(map[string]*websocket.Conn),
		broadcast:  make(chan []byte),
		register:   make(chan *WebSocketPacket),
		unregister: make(chan *WebSocketPacket),
	}
}

func (h *webSocketHub) New(c *echo.Context) error {
	conn, err := NewWebSocketUpgrader().Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	// 註冊 Client
	h.register <- &WebSocketPacket{ID: "q22312", Conn: conn}

	// 處理離開時的清除機制與關閉連線
	defer func() {
		h.unregister <- &WebSocketPacket{ID: "q22312", Conn: conn}
		conn.Close()
	}()

	// 持續讀取該 Client 發送過來的訊息
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			// 斷開連線或發生錯誤時跳出迴圈
			break
		}
		// 收到任何訊息，丟進廣播 channel 發給所有人
		h.broadcast <- message
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
