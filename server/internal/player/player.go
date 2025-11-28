package player

import (
	"math/rand"
	"sync"
	"time"
)

// Player 玩家结构
type Player struct {
	ID          string
	Username    string
	Color       Color
	Position    Vector3
	Rotation    Quaternion
	VehicleType string      // 车辆类型：ToyCar, Bus, GoKart, Bulldozer, Crane, Truck
	Conn        interface{} // WebSocket连接
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

// Color 颜色
type Color struct {
	R float32
	G float32
	B float32
}

// Manager 玩家管理器
type Manager struct {
	players map[string]*Player
	mu      sync.RWMutex
}

// NewManager 创建玩家管理器
func NewManager() *Manager {
	rand.Seed(time.Now().UnixNano())
	return &Manager{
		players: make(map[string]*Player),
	}
}

// AddPlayer 添加玩家
func (m *Manager) AddPlayer(player *Player) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.players[player.ID] = player
}

// RemovePlayer 移除玩家
func (m *Manager) RemovePlayer(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.players, id)
}

// GetPlayer 获取玩家
func (m *Manager) GetPlayer(id string) (*Player, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	player, ok := m.players[id]
	return player, ok
}

// GetAllPlayers 获取所有玩家
func (m *Manager) GetAllPlayers() []*Player {
	m.mu.RLock()
	defer m.mu.RUnlock()
	players := make([]*Player, 0, len(m.players))
	for _, p := range m.players {
		players = append(players, p)
	}
	return players
}

// UpdatePlayerPosition 更新玩家位置
func (m *Manager) UpdatePlayerPosition(id string, pos Vector3, rot Quaternion) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if player, ok := m.players[id]; ok {
		player.Position = pos
		player.Rotation = rot
		return true
	}
	return false
}

// GenerateUsername 生成随机用户名
func GenerateUsername() string {
	adjectives := []string{
		"Fast", "Quick", "Speedy", "Turbo", "Racing",
		"Cool", "Super", "Mega", "Ultra", "Epic",
		"Wild", "Crazy", "Swift", "Rapid", "Lightning",
	}
	nouns := []string{
		"Racer", "Driver", "Pilot", "Rider", "Cruiser",
		"Drifter", "Speedster", "Chaser", "Runner", "Zoomer",
		"Warrior", "Champion", "Master", "Hero", "Legend",
	}

	adj := adjectives[rand.Intn(len(adjectives))]
	noun := nouns[rand.Intn(len(nouns))]
	num := rand.Intn(1000)

	return adj + noun + string(rune('0'+num/100)) + string(rune('0'+(num/10)%10)) + string(rune('0'+num%10))
}

// GenerateRandomColor 生成随机颜色
func GenerateRandomColor() Color {
	return Color{
		R: rand.Float32()*0.7 + 0.3, // 0.3-1.0
		G: rand.Float32()*0.7 + 0.3,
		B: rand.Float32()*0.7 + 0.3,
	}
}

// GenerateSpawnPosition 生成出生位置（场地中央附近）
func GenerateSpawnPosition() Vector3 {
	// 在中央附近随机生成位置
	offsetX := (rand.Float32() - 0.5) * 10 // -5 到 5
	offsetZ := (rand.Float32() - 0.5) * 10
	return Vector3{
		X: offsetX,
		Y: 1.2, // 车辆高度
		Z: offsetZ,
	}
}

// GenerateRandomVehicleType 生成随机车辆类型
func GenerateRandomVehicleType() string {
	vehicleTypes := []string{
		"ToyCar",
		"Bus",
		"GoKart",
		"Bulldozer",
		"Crane",
		"Truck",
	}
	return vehicleTypes[rand.Intn(len(vehicleTypes))]
}
