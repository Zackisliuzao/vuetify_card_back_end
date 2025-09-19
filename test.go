package main

import "fmt"

// 1. 定义一个接口（抽象的厨师）
type Cook interface {
	CookMeal() string
}

// 2. 实现一个真正的厨师
type RealCook struct{}

func (c *RealCook) CookMeal() string {
	return "🍖 红烧肉"
}

// 3. 实现一个假的厨师（用于测试或特殊场景）
type MockCook struct{}

func (m *MockCook) CookMeal() string {
	return "🥪 测试餐（假数据）"
}

// 4. 经理，依赖一个 Cook（通过接口，不是具体实现！）
type Manager struct {
	cook Cook // 只依赖接口，不依赖具体是谁
}

// 5. 通过构造函数注入 Cook（重点！这就是依赖注入！）
func NewManager(cook Cook) *Manager {
	return &Manager{cook: cook}
}

// 6. 经理调用厨师做饭
func (m *Manager) ServeMeal() string {
	return m.cook.CookMeal()
}

func main() {
	// 场景1：用真厨师
	realCook := &RealCook{}
	manager := NewManager(realCook)
	fmt.Println(manager.ServeMeal()) // 输出：🍖 红烧肉

	// 场景2：用假厨师（比如测试时）
	mockCook := &MockCook{}
	manager2 := NewManager(mockCook)
	fmt.Println(manager2.ServeMeal()) // 输出：🥪 测试餐（假数据）
}
