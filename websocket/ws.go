package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/gorilla/websocket"
)

// 设备数据结构
type DeviceData struct {
	DeviceID    string  `json:"deviceId"`
	Status      string  `json:"status"`
	Temperature float64 `json:"temperature"`
}

// 客户端发送的消息结构（扩展）
type ClientMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"` // 延迟解析
}

// 降低温度的数据结构
type DecreaseTempData struct {
	DeviceID string  `json:"deviceId"`
	Amount   float64 `json:"amount"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	log.Println("WebSocket connected")

	// 存储所有设备的温度数据
	deviceTemperatures := make(map[string]DeviceData)
	// 初始化一些默认设备
	deviceTemperatures["笔记本"] = DeviceData{DeviceID: "笔记本", Status: "在线", Temperature: 35.0}
	deviceTemperatures["手机"] = DeviceData{DeviceID: "手机", Status: "在线", Temperature: 36.0}
	deviceTemperatures["耳机"] = DeviceData{DeviceID: "耳机", Status: "在线", Temperature: 37.0}

	// 收消息 goroutine
	go func() {
		for {
			_, messageR, err := conn.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				// 这里如果退出，主 goroutine也要退出
				conn.Close()
				return
			}

			var msg ClientMessage
			if err := json.Unmarshal(messageR, &msg); err != nil {
				log.Println("JSON unmarshal error:", err)
				continue
			}
			fmt.Println("Received message:", string(messageR))

			switch msg.Type {
			case "decrease_temperature":
				var data DecreaseTempData
				if err := json.Unmarshal(msg.Data, &data); err != nil {
					log.Println("Failed to parse decrease_temperature data:", err)
					continue
				}
				log.Printf("收到降温指令: 设备 %s 温度减少 %.1f", data.DeviceID, data.Amount)
				if device, ok := deviceTemperatures[data.DeviceID]; ok {
					device.Temperature -= data.Amount
					deviceTemperatures[data.DeviceID] = device
				} else {
					log.Printf("设备 %s 不存在，无法降温", data.DeviceID)
				}

			default:
				log.Printf("未知消息类型: %s", msg.Type)
			}
		}
	}()

	// 主循环：每秒发送
	// ticker := time.NewTicker(1 * time.Second)
	// defer ticker.Stop()

	for {
		// 将map中的设备数据转换为数组
		var keys []string
		for k := range deviceTemperatures {
			keys = append(keys, k)
		}
		// 排序
		sort.Strings(keys)

		// 遍历排序好的 keys
		var allDevices []DeviceData
		for _, k := range keys {
			device := deviceTemperatures[k]
			device.Temperature += 1        // 模拟温度上升
			deviceTemperatures[k] = device // 写回去
			allDevices = append(allDevices, device)
		}

		message, err := json.Marshal(allDevices)
		if err != nil {
			log.Println("JSON marshal error:", err)
			break
		}

		err = conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("Write error:", err)
			break
		}

		time.Sleep(1 * time.Second)
	}
}

func main() {
	http.HandleFunc("/ws", wsHandler)

	addr := ":8090" // 👈 改这里
	log.Println("Server started at", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
