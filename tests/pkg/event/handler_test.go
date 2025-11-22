package event_test

import (
	"testing"

	pkgEvent "github.com/mautops/approval-kit/pkg/event"
	internalEvent "github.com/mautops/approval-kit/internal/event"
)

// TestPkgEventHandlerInterface 验证 pkg/event 包中的 EventHandler 接口定义
func TestPkgEventHandlerInterface(t *testing.T) {
	// 验证接口类型存在
	var handler pkgEvent.EventHandler
	if handler != nil {
		_ = handler
	}
}

// TestPkgEventHandlerCompatibility 验证 pkg/event.EventHandler 与 internal/event.EventHandler 兼容
func TestPkgEventHandlerCompatibility(t *testing.T) {
	// 验证 pkg 接口可以被 internal 实现满足
	var _ pkgEvent.EventHandler = (*internalEventHandlerAdapter)(nil)
}

// internalEventHandlerAdapter 用于测试接口兼容性的适配器
type internalEventHandlerAdapter struct {
	impl internalEvent.EventHandler
}

func (a *internalEventHandlerAdapter) Handle(evt *pkgEvent.Event) error {
	internalEvt := pkgEvent.EventToInternal(evt)
	return a.impl.Handle(internalEvt)
}

