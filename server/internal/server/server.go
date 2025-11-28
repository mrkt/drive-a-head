package server

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"server/internal/player"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源（生产环境需要更严格的检查）
	},
}

// Client WebSocket客户端
type Client struct {
	ID       string
	Conn     *websocket.Conn
	Server   *Server
	SendChan chan []byte
	LastSeen int64 // 最后活跃时间（Unix时间戳）
}

// Server 游戏服务器
type Server struct {
	clients       map[string]*Client
	playerManager *player.Manager
	worldManager  *SimpleWorldManager
	broadcast     chan *BroadcastMessage
	register      chan *Client
	unregister    chan *Client
	mu            sync.RWMutex
}

// WorldManager 世界管理器接口
type WorldManager interface {
	GetAllObjects() []WorldObject
}

// WorldObject 世界对象
type WorldObject struct {
	ID       string
	Type     string
	Position Vector3
	Rotation Quaternion
	Scale    Vector3
}

// Vector3 三维向量
type Vector3 struct {
	X float32
	Y float32
	Z float32
}

// Quaternion 四元数
type Quaternion struct {
	X float32
	Y float32
	Z float32
	W float32
}

// BroadcastMessage 广播消息
type BroadcastMessage struct {
	Data       []byte
	ExcludeID  string // 排除的客户端ID
	IncludeIDs []string // 只发送给指定的客户端ID列表（如果为空则发送给所有人）
}

// NewServer 创建服务器
func NewServer() *Server {
	return &Server{
		clients:       make(map[string]*Client),
		playerManager: player.NewManager(),
		worldManager:  NewSimpleWorldManager(),
		broadcast:     make(chan *BroadcastMessage, 256),
		register:      make(chan *Client),
		unregister:    make(chan *Client),
	}
}

// SimpleWorldManager 简单的世界管理器实现
type SimpleWorldManager struct {
	objects []WorldObject
}

// NewSimpleWorldManager 创建简单世界管理器
func NewSimpleWorldManager() *SimpleWorldManager {
	// 场景已经在 Unity 中预设，不需要服务器同步静态对象
	return &SimpleWorldManager{
		objects: []WorldObject{},
	}
}

// GetAllObjects 获取所有对象
func (w *SimpleWorldManager) GetAllObjects() []WorldObject {
	return w.objects
}

// Run 运行服务器
func (s *Server) Run() {
	for {
		select {
		case client := <-s.register:
			s.mu.Lock()
			s.clients[client.ID] = client
			s.mu.Unlock()
			log.Printf("Client registered: %s", client.ID)

		case client := <-s.unregister:
			s.mu.Lock()
			if _, ok := s.clients[client.ID]; ok {
				delete(s.clients, client.ID)
				close(client.SendChan)
				s.playerManager.RemovePlayer(client.ID)
				log.Printf("Client unregistered: %s", client.ID)

				// 通知其他玩家该玩家离开
				s.broadcastPlayerRemove(client.ID)
			}
			s.mu.Unlock()

		case message := <-s.broadcast:
			s.mu.RLock()
			for id, client := range s.clients {
				// 如果指定了排除ID，跳过该客户端
				if message.ExcludeID != "" && id == message.ExcludeID {
					continue
				}

				// 如果指定了包含ID列表，只发送给列表中的客户端
				if len(message.IncludeIDs) > 0 {
					found := false
					for _, includeID := range message.IncludeIDs {
						if id == includeID {
							found = true
							break
						}
					}
					if !found {
						continue
					}
				}

				select {
				case client.SendChan <- message.Data:
				default:
					close(client.SendChan)
					delete(s.clients, id)
				}
			}
			s.mu.RUnlock()
		}
	}
}

// HandleWebSocket 处理WebSocket连接
func (s *Server) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	log.Printf("WebSocket connection request from %s", r.RemoteAddr)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	// 生成唯一ID
	clientID := generateClientID()
	client := &Client{
		ID:       clientID,
		Conn:     conn,
		Server:   s,
		SendChan: make(chan []byte, 256),
		LastSeen: time.Now().Unix(),
	}

	log.Printf("New WebSocket client connected: %s", clientID)

	s.register <- client

	// 启动发送和接收协程
	go client.writePump()
	go client.readPump()
}

// writePump 发送消息到客户端
func (c *Client) writePump() {
	defer func() {
		c.Conn.Close()
	}()

	log.Printf("Client %s writePump started", c.ID)

	for message := range c.SendChan {
		err := c.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("Write error to %s: %v", c.ID, err)
			return
		}
		log.Printf("Sent %d bytes to client %s", len(message), c.ID)
	}
}

// readPump 从客户端接收消息
func (c *Client) readPump() {
	defer func() {
		c.Server.unregister <- c
		c.Conn.Close()
	}()

	log.Printf("Client %s readPump started", c.ID)

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Read error from %s: %v", c.ID, err)
			}
			break
		}

		log.Printf("Received %d bytes from client %s", len(message), c.ID)
		c.handleMessage(message)
	}
}

// generateClientID 生成客户端ID（使用UUID）
func generateClientID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// SendMessage 发送消息给客户端
func (c *Client) SendMessage(data []byte) {
	select {
	case c.SendChan <- data:
	default:
		log.Printf("Client %s send channel full", c.ID)
	}
}

// Broadcast 广播消息
func (s *Server) Broadcast(data []byte, excludeID string) {
	s.broadcast <- &BroadcastMessage{
		Data:      data,
		ExcludeID: excludeID,
	}
}

// BroadcastToClients 广播消息给指定客户端
func (s *Server) BroadcastToClients(data []byte, clientIDs []string) {
	s.broadcast <- &BroadcastMessage{
		Data:       data,
		IncludeIDs: clientIDs,
	}
}

// MarshalProtoMessage 序列化protobuf消息
func MarshalProtoMessage(msg proto.Message) ([]byte, error) {
	return proto.Marshal(msg)
}

// UnmarshalProtoMessage 反序列化protobuf消息
func UnmarshalProtoMessage(data []byte, msg proto.Message) error {
	return proto.Unmarshal(data, msg)
}
