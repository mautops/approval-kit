package template_test

import (
	"testing"

	"github.com/mautops/approval-kit/internal/template"
)

// TestTemplateManagerInterface 验证 TemplateManager 接口定义
func TestTemplateManagerInterface(t *testing.T) {
	// 验证接口类型存在
	var mgr template.TemplateManager
	if mgr != nil {
		_ = mgr
	}
}

// TestTemplateManagerMethods 验证 TemplateManager 接口方法签名
func TestTemplateManagerMethods(t *testing.T) {
	// 验证接口包含所有必需的方法
	// 通过编译时检查,如果方法不存在会编译失败
	var _ template.TemplateManager = (*templateManagerImpl)(nil)
}

// templateManagerImpl 用于测试接口方法签名的实现
type templateManagerImpl struct{}

func (m *templateManagerImpl) Create(tpl *template.Template) error {
	return nil
}

func (m *templateManagerImpl) Update(id string, tpl *template.Template) error {
	return nil
}

func (m *templateManagerImpl) Get(id string, version int) (*template.Template, error) {
	return nil, nil
}

func (m *templateManagerImpl) Delete(id string) error {
	return nil
}

func (m *templateManagerImpl) ListVersions(id string) ([]int, error) {
	return nil, nil
}

