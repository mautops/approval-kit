package template_test

import (
	"testing"

	"github.com/mautops/approval-kit/internal/template"
)

// TestTemplateManagerDelete 测试模板删除功能
func TestTemplateManagerDelete(t *testing.T) {
	mgr := template.NewTemplateManager()

	// 创建模板
	tpl := createTestTemplate()
	err := mgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// 删除模板
	err = mgr.Delete(tpl.ID)
	if err != nil {
		t.Fatalf("Delete() failed: %v", err)
	}

	// 验证模板已删除
	_, err = mgr.Get(tpl.ID, 0)
	if err == nil {
		t.Error("Get() should fail after Delete()")
	}
}

// TestTemplateManagerDeleteAllVersions 测试删除所有版本
func TestTemplateManagerDeleteAllVersions(t *testing.T) {
	mgr := template.NewTemplateManager()

	// 创建多个版本
	tpl1 := createTestTemplate()
	tpl1.Version = 1
	err := mgr.Create(tpl1)
	if err != nil {
		t.Fatalf("Create() failed for version 1: %v", err)
	}

	tpl2 := createTestTemplate()
	tpl2.Version = 2
	err = mgr.Create(tpl2)
	if err != nil {
		t.Fatalf("Create() failed for version 2: %v", err)
	}

	// 删除模板(应该删除所有版本)
	err = mgr.Delete(tpl1.ID)
	if err != nil {
		t.Fatalf("Delete() failed: %v", err)
	}

	// 验证所有版本都已删除
	_, err = mgr.Get(tpl1.ID, 1)
	if err == nil {
		t.Error("Get() should fail for version 1 after Delete()")
	}

	_, err = mgr.Get(tpl1.ID, 2)
	if err == nil {
		t.Error("Get() should fail for version 2 after Delete()")
	}

	// 验证版本列表为空
	versions, err := mgr.ListVersions(tpl1.ID)
	if err == nil && len(versions) > 0 {
		t.Error("ListVersions() should return empty list after Delete()")
	}
}

// TestTemplateManagerDeleteNotFound 测试删除不存在的模板
func TestTemplateManagerDeleteNotFound(t *testing.T) {
	mgr := template.NewTemplateManager()

	// 尝试删除不存在的模板
	err := mgr.Delete("non-existent")
	if err == nil {
		t.Error("Delete() should fail when template does not exist")
	}
}

// TestTemplateManagerDeleteAfterUpdate 测试更新后删除
func TestTemplateManagerDeleteAfterUpdate(t *testing.T) {
	mgr := template.NewTemplateManager()

	// 创建模板
	tpl := createTestTemplate()
	err := mgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// 更新模板
	updatedTpl := createTestTemplate()
	updatedTpl.Name = "Updated Template"
	err = mgr.Update(tpl.ID, updatedTpl)
	if err != nil {
		t.Fatalf("Update() failed: %v", err)
	}

	// 删除模板
	err = mgr.Delete(tpl.ID)
	if err != nil {
		t.Fatalf("Delete() failed: %v", err)
	}

	// 验证所有版本都已删除
	_, err = mgr.Get(tpl.ID, 1)
	if err == nil {
		t.Error("Get() should fail for version 1 after Delete()")
	}

	_, err = mgr.Get(tpl.ID, 2)
	if err == nil {
		t.Error("Get() should fail for version 2 after Delete()")
	}
}

