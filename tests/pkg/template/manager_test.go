package template_test

import (
	"testing"

	pkgTemplate "github.com/mautops/approval-kit/pkg/template"
	internalTemplate "github.com/mautops/approval-kit/internal/template"
)

// TestPkgTemplateManagerInterface 验证 pkg/template 包中的 TemplateManager 接口定义
func TestPkgTemplateManagerInterface(t *testing.T) {
	// 验证接口类型存在
	var mgr pkgTemplate.TemplateManager
	if mgr != nil {
		_ = mgr
	}
}

// TestPkgTemplateManagerCompatibility 验证 pkg/template.TemplateManager 与 internal/template.TemplateManager 兼容
func TestPkgTemplateManagerCompatibility(t *testing.T) {
	// 验证 pkg 接口可以被 internal 实现满足
	var _ pkgTemplate.TemplateManager = (*internalTemplateManagerAdapter)(nil)
}

// internalTemplateManagerAdapter 用于测试接口兼容性的适配器
type internalTemplateManagerAdapter struct {
	impl internalTemplate.TemplateManager
}

func (a *internalTemplateManagerAdapter) Create(template *pkgTemplate.Template) error {
	internalTpl := pkgTemplate.TemplateToInternal(template)
	return a.impl.Create(internalTpl)
}

func (a *internalTemplateManagerAdapter) Update(id string, template *pkgTemplate.Template) error {
	internalTpl := pkgTemplate.TemplateToInternal(template)
	return a.impl.Update(id, internalTpl)
}

func (a *internalTemplateManagerAdapter) Get(id string, version int) (*pkgTemplate.Template, error) {
	internalTpl, err := a.impl.Get(id, version)
	if err != nil {
		return nil, err
	}
	return pkgTemplate.TemplateFromInternal(internalTpl), nil
}

func (a *internalTemplateManagerAdapter) Delete(id string) error {
	return a.impl.Delete(id)
}

func (a *internalTemplateManagerAdapter) ListVersions(id string) ([]int, error) {
	return a.impl.ListVersions(id)
}

