package protocol

import (
	"encoding/binary"
	"encoding/json"
)

// EncodeMessage 编码消息（使用简单的JSON格式，便于Unity解析）
func EncodeMessage(msgType MessageType, payload interface{}) ([]byte, error) {
	msg := map[string]interface{}{
		"type":    int32(msgType),
		"payload": payload,
	}
	return json.Marshal(msg)
}

// DecodeMessage 解码消息
func DecodeMessage(data []byte) (*GameMessage, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	msgType := MessageType(int32(raw["type"].(float64)))

	msg := &GameMessage{
		Type: msgType,
	}

	// 根据消息类型解析payload
	if payload, ok := raw["payload"]; ok && payload != nil {
		payloadBytes, _ := json.Marshal(payload)

		switch msgType {
		case MessageTypeJoin:
			// 加入请求没有payload
			msg.Payload = nil

		case MessageTypePlayerUpdate:
			var req PlayerUpdateRequest
			if err := json.Unmarshal(payloadBytes, &req); err == nil {
				msg.Payload = &req
			}
		}
	}

	return msg, nil
}

// EncodeJoinResponse 编码加入响应
func EncodeJoinResponse(resp *JoinResponse) ([]byte, error) {
	payload := map[string]interface{}{
		"player_id": resp.PlayerID,
		"username":  resp.Username,
		"color": map[string]float32{
			"r": resp.Color.R,
			"g": resp.Color.G,
			"b": resp.Color.B,
		},
		"spawn_position": map[string]float32{
			"x": resp.SpawnPosition.X,
			"y": resp.SpawnPosition.Y,
			"z": resp.SpawnPosition.Z,
		},
		"vehicle_type": resp.VehicleType,
	}
	return EncodeMessage(MessageTypeJoin, payload)
}

// EncodePlayerUpdateBroadcast 编码玩家更新广播
func EncodePlayerUpdateBroadcast(broadcast *PlayerUpdateBroadcast) ([]byte, error) {
	payload := map[string]interface{}{
		"player_id": broadcast.PlayerID,
		"position": map[string]float32{
			"x": broadcast.Position.X,
			"y": broadcast.Position.Y,
			"z": broadcast.Position.Z,
		},
		"rotation": map[string]float32{
			"x": broadcast.Rotation.X,
			"y": broadcast.Rotation.Y,
			"z": broadcast.Rotation.Z,
			"w": broadcast.Rotation.W,
		},
	}
	return EncodeMessage(MessageTypePlayerUpdate, payload)
}

// EncodeSpawnPlayerNotification 编码生成玩家通知
func EncodeSpawnPlayerNotification(notification *SpawnPlayerNotification) ([]byte, error) {
	p := notification.Player
	payload := map[string]interface{}{
		"player": map[string]interface{}{
			"id":       p.ID,
			"username": p.Username,
			"color": map[string]float32{
				"r": p.Color.R,
				"g": p.Color.G,
				"b": p.Color.B,
			},
			"position": map[string]float32{
				"x": p.Position.X,
				"y": p.Position.Y,
				"z": p.Position.Z,
			},
			"rotation": map[string]float32{
				"x": p.Rotation.X,
				"y": p.Rotation.Y,
				"z": p.Rotation.Z,
				"w": p.Rotation.W,
			},
			"vehicle_type": p.VehicleType,
		},
	}
	return EncodeMessage(MessageTypeSpawnPlayer, payload)
}

// EncodeRemovePlayerNotification 编码移除玩家通知
func EncodeRemovePlayerNotification(notification *RemovePlayerNotification) ([]byte, error) {
	payload := map[string]interface{}{
		"player_id": notification.PlayerID,
	}
	return EncodeMessage(MessageTypeRemovePlayer, payload)
}

// EncodePlayerListResponse 编码玩家列表响应
func EncodePlayerListResponse(resp *PlayerListResponse) ([]byte, error) {
	players := make([]map[string]interface{}, 0, len(resp.Players))

	for _, p := range resp.Players {
		players = append(players, map[string]interface{}{
			"id":       p.ID,
			"username": p.Username,
			"color": map[string]float32{
				"r": p.Color.R,
				"g": p.Color.G,
				"b": p.Color.B,
			},
			"position": map[string]float32{
				"x": p.Position.X,
				"y": p.Position.Y,
				"z": p.Position.Z,
			},
			"rotation": map[string]float32{
				"x": p.Rotation.X,
				"y": p.Rotation.Y,
				"z": p.Rotation.Z,
				"w": p.Rotation.W,
			},
			"vehicle_type": p.VehicleType,
		})
	}

	payload := map[string]interface{}{
		"players": players,
	}
	return EncodeMessage(MessageTypePlayerList, payload)
}

// EncodeWorldStateResponse 编码世界状态响应
func EncodeWorldStateResponse(resp *WorldStateResponse) ([]byte, error) {
	objects := make([]map[string]interface{}, 0, len(resp.Objects))

	for _, obj := range resp.Objects {
		objects = append(objects, map[string]interface{}{
			"id":   obj.ID,
			"type": obj.Type,
			"position": map[string]float32{
				"x": obj.Position.X,
				"y": obj.Position.Y,
				"z": obj.Position.Z,
			},
			"rotation": map[string]float32{
				"x": obj.Rotation.X,
				"y": obj.Rotation.Y,
				"z": obj.Rotation.Z,
				"w": obj.Rotation.W,
			},
			"scale": map[string]float32{
				"x": obj.Scale.X,
				"y": obj.Scale.Y,
				"z": obj.Scale.Z,
			},
		})
	}

	payload := map[string]interface{}{
		"objects": objects,
	}
	return EncodeMessage(MessageTypeWorldState, payload)
}

// EncodeBinary 编码为二进制（带长度前缀）
func EncodeBinary(data []byte) []byte {
	length := uint32(len(data))
	result := make([]byte, 4+length)
	binary.BigEndian.PutUint32(result[0:4], length)
	copy(result[4:], data)
	return result
}
