package event

// EventHandler 事件处理器接口
// 用于处理审批流程中的事件通知
// 与 internal/event.EventHandler 接口定义完全一致,但位于 pkg 目录,可以被外部导入
type EventHandler interface {
	// Handle 处理事件
	// evt: 事件对象
	// 返回: 错误信息
	Handle(evt *Event) error
}

