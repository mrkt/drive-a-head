package server

import (
	"log"
	"server/internal/protocol"
)

// handleMessage 处理接收到的消息（完整实现）
func (c *Client) handleMessage(data []byte) {
	// 解码消息
	msg, err := protocol.DecodeMessage(data)
	if err != nil {
		log.Printf("Failed to decode message from client %s: %v", c.ID, err)
		return
	}

	switch msg.Type {
	case protocol.MessageTypeJoin:
		c.handleJoin()

	case protocol.MessageTypePlayerUpdate:
		if req, ok := msg.Payload.(*protocol.PlayerUpdateRequest); ok {
			c.handlePlayerUpdate(req)
		}

	default:
		log.Printf("Unknown message type from client %s: %d", c.ID, msg.Type)
	}
}

// handleJoin 处理玩家加入
func (c *Client) handleJoin() {
	// 处理加入请求
	resp := protocol.HandleJoinRequest(c.ID, c.Server.playerManager)

	// 发送加入响应给当前玩家
	respData, err := protocol.EncodeJoinResponse(resp)
	if err != nil {
		log.Printf("Failed to encode join response: %v", err)
		return
	}
	c.SendMessage(respData)

	// 发送世界状态给新玩家
	worldObjects := c.Server.worldManager.GetAllObjects()
	protoObjects := make([]protocol.GameObject, len(worldObjects))
	for i, obj := range worldObjects {
		protoObjects[i] = protocol.GameObject{
			ID:   obj.ID,
			Type: obj.Type,
			Position: protocol.Vector3{
				X: obj.Position.X,
				Y: obj.Position.Y,
				Z: obj.Position.Z,
			},
			Rotation: protocol.Quaternion{
				X: obj.Rotation.X,
				Y: obj.Rotation.Y,
				Z: obj.Rotation.Z,
				W: obj.Rotation.W,
			},
			Scale: protocol.Vector3{
				X: obj.Scale.X,
				Y: obj.Scale.Y,
				Z: obj.Scale.Z,
			},
		}
	}
	worldState := &protocol.WorldStateResponse{
		Objects: protoObjects,
	}
	worldData, err := protocol.EncodeWorldStateResponse(worldState)
	if err != nil {
		log.Printf("Failed to encode world state: %v", err)
	} else {
		c.SendMessage(worldData)
		log.Printf("Sent world state with %d objects to %s", len(protoObjects), c.ID)
	}

	// 获取现有玩家列表并发送给新玩家
	playerList := protocol.GetPlayerList(c.Server.playerManager, c.ID)
	listData, err := protocol.EncodePlayerListResponse(playerList)
	if err != nil {
		log.Printf("Failed to encode player list: %v", err)
	} else {
		c.SendMessage(listData)
	}

	// 通知其他玩家有新玩家加入
	player, _ := c.Server.playerManager.GetPlayer(c.ID)
	spawnNotification := &protocol.SpawnPlayerNotification{
		Player: player,
	}
	spawnData, err := protocol.EncodeSpawnPlayerNotification(spawnNotification)
	if err != nil {
		log.Printf("Failed to encode spawn notification: %v", err)
	} else {
		c.Server.Broadcast(spawnData, c.ID)
	}

	log.Printf("Player %s joined successfully", resp.Username)
}

// handlePlayerUpdate 处理玩家位置更新
func (c *Client) handlePlayerUpdate(req *protocol.PlayerUpdateRequest) {
	// 更新玩家位置
	broadcast := protocol.HandlePlayerUpdate(c.ID, req, c.Server.playerManager)
	if broadcast == nil {
		return
	}

	// 广播给其他玩家
	broadcastData, err := protocol.EncodePlayerUpdateBroadcast(broadcast)
	if err != nil {
		log.Printf("Failed to encode player update broadcast: %v", err)
		return
	}

	c.Server.Broadcast(broadcastData, c.ID)
}

// broadcastPlayerRemove 广播玩家离开（完整实现）
func (s *Server) broadcastPlayerRemove(playerID string) {
	notification := &protocol.RemovePlayerNotification{
		PlayerID: playerID,
	}

	data, err := protocol.EncodeRemovePlayerNotification(notification)
	if err != nil {
		log.Printf("Failed to encode remove player notification: %v", err)
		return
	}

	s.Broadcast(data, playerID)
	log.Printf("Broadcasted player remove: %s", playerID)
}
