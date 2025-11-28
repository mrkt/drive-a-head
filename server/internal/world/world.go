package world

import (
	"sync"
)

// GameObject 场景中的游戏对象
type GameObject struct {
	ID       string
	Type     string // "obstacle", "prop", "ground"
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

// World 游戏世界状态
type World struct {
	objects map[string]*GameObject
	mu      sync.RWMutex
}

// NewWorld 创建游戏世界
func NewWorld() *World {
	w := &World{
		objects: make(map[string]*GameObject),
	}
	w.initializeDefaultWorld()
	return w
}

// initializeDefaultWorld 初始化默认世界（地面和一些障碍物）
func (w *World) initializeDefaultWorld() {
	// 添加地面
	w.objects["ground"] = &GameObject{
		ID:   "ground",
		Type: "ground",
		Position: Vector3{X: 0, Y: 0, Z: 0},
		Rotation: Quaternion{X: 0, Y: 0, Z: 0, W: 1},
		Scale:    Vector3{X: 50, Y: 1, Z: 50},
	}

	// 添加一些障碍物
	obstacles := []struct {
		id  string
		pos Vector3
	}{
		{"obstacle1", Vector3{X: 10, Y: 1, Z: 10}},
		{"obstacle2", Vector3{X: -10, Y: 1, Z: 10}},
		{"obstacle3", Vector3{X: 10, Y: 1, Z: -10}},
		{"obstacle4", Vector3{X: -10, Y: 1, Z: -10}},
	}

	for _, obs := range obstacles {
		w.objects[obs.id] = &GameObject{
			ID:       obs.id,
			Type:     "obstacle",
			Position: obs.pos,
			Rotation: Quaternion{X: 0, Y: 0, Z: 0, W: 1},
			Scale:    Vector3{X: 2, Y: 2, Z: 2},
		}
	}
}

// GetAllObjects 获取所有游戏对象
func (w *World) GetAllObjects() []*GameObject {
	w.mu.RLock()
	defer w.mu.RUnlock()

	objects := make([]*GameObject, 0, len(w.objects))
	for _, obj := range w.objects {
		objects = append(objects, obj)
	}
	return objects
}

// AddObject 添加游戏对象
func (w *World) AddObject(obj *GameObject) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.objects[obj.ID] = obj
}

// RemoveObject 移除游戏对象
func (w *World) RemoveObject(id string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	delete(w.objects, id)
}
