package protocol

import (
	"log"
	"server/internal/player"
)

// MessageType 消息类型
type MessageType int32

const (
	MessageTypeUnknown       MessageType = 0
	MessageTypeJoin          MessageType = 1
	MessageTypeLeave         MessageType = 2
	MessageTypePlayerUpdate  MessageType = 3
	MessageTypeSpawnPlayer   MessageType = 4
	MessageTypeRemovePlayer  MessageType = 5
	MessageTypePlayerList    MessageType = 6
	MessageTypeWorldState    MessageType = 7
)

// GameMessage 游戏消息（简化版，不依赖protobuf生成的代码）
type GameMessage struct {
	Type    MessageType
	Payload interface{}
}

// JoinResponse 加入响应
type JoinResponse struct {
	PlayerID      string
	Username      string
	Color         player.Color
	SpawnPosition player.Vector3
	VehicleType   string
}

// PlayerUpdateRequest 玩家更新请求
type PlayerUpdateRequest struct {
	Position player.Vector3
	Rotation player.Quaternion
}

// PlayerUpdateBroadcast 玩家更新广播
type PlayerUpdateBroadcast struct {
	PlayerID string
	Position player.Vector3
	Rotation player.Quaternion
}

// SpawnPlayerNotification 生成玩家通知
type SpawnPlayerNotification struct {
	Player *player.Player
}

// RemovePlayerNotification 移除玩家通知
type RemovePlayerNotification struct {
	PlayerID string
}

// PlayerListResponse 玩家列表响应
type PlayerListResponse struct {
	Players []*player.Player
}

// GameObject 场景对象
type GameObject struct {
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

// WorldStateResponse 世界状态响应
type WorldStateResponse struct {
	Objects []GameObject
}

// HandleJoinRequest 处理加入请求
func HandleJoinRequest(clientID string, pm *player.Manager) *JoinResponse {
	// 生成玩家数据
	username := player.GenerateUsername()
	color := player.GenerateRandomColor()
	spawnPos := player.GenerateSpawnPosition()
	vehicleType := player.GenerateRandomVehicleType()

	// 创建玩家
	newPlayer := &player.Player{
		ID:          clientID,
		Username:    username,
		Color:       color,
		Position:    spawnPos,
		Rotation:    player.Quaternion{X: 0, Y: 0, Z: 0, W: 1},
		VehicleType: vehicleType,
	}

	// 添加到玩家管理器
	pm.AddPlayer(newPlayer)

	log.Printf("Player joined: %s (%s) with vehicle: %s", username, clientID, vehicleType)

	return &JoinResponse{
		PlayerID:      clientID,
		Username:      username,
		Color:         color,
		SpawnPosition: spawnPos,
		VehicleType:   vehicleType,
	}
}

// HandlePlayerUpdate 处理玩家更新
func HandlePlayerUpdate(clientID string, req *PlayerUpdateRequest, pm *player.Manager) *PlayerUpdateBroadcast {
	// 更新玩家位置
	if pm.UpdatePlayerPosition(clientID, req.Position, req.Rotation) {
		return &PlayerUpdateBroadcast{
			PlayerID: clientID,
			Position: req.Position,
			Rotation: req.Rotation,
		}
	}
	return nil
}

// GetPlayerList 获取玩家列表
func GetPlayerList(pm *player.Manager, excludeID string) *PlayerListResponse {
	allPlayers := pm.GetAllPlayers()
	players := make([]*player.Player, 0)

	for _, p := range allPlayers {
		if p.ID != excludeID {
			players = append(players, p)
		}
	}

	return &PlayerListResponse{
		Players: players,
	}
}
