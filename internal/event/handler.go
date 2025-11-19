package event

// EventHandler 事件处理器接口
// 用于处理审批流程中的事件通知
type EventHandler interface {
	// Handle 处理事件
	// evt: 事件对象
	// 返回: 错误信息
	Handle(evt *Event) error
}

