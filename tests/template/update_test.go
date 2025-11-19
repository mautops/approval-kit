package template_test

import (
	"testing"
	"time"

	"github.com/mautops/approval-kit/internal/template"
)

// TestTemplateManagerUpdate 测试模板更新功能
func TestTemplateManagerUpdate(t *testing.T) {
	mgr := template.NewTemplateManager()

	// 先创建模板
	tpl := createTestTemplate()
	tpl.Version = 1
	err := mgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// 更新模板
	updatedTpl := createTestTemplate()
	updatedTpl.Name = "Updated Template"
	updatedTpl.Description = "Updated description"
	updatedTpl.Version = 2 // 新版本

	err = mgr.Update(tpl.ID, updatedTpl)
	if err != nil {
		t.Fatalf("Update() failed: %v", err)
	}

	// 验证原版本仍然存在
	original, err := mgr.Get(tpl.ID, 1)
	if err != nil {
		t.Fatalf("Get() failed for original version: %v", err)
	}
	if original.Name != tpl.Name {
		t.Errorf("Original template Name = %q, want %q", original.Name, tpl.Name)
	}

	// 验证新版本已创建
	updated, err := mgr.Get(tpl.ID, 2)
	if err != nil {
		t.Fatalf("Get() failed for updated version: %v", err)
	}
	if updated.Name != "Updated Template" {
		t.Errorf("Updated template Name = %q, want %q", updated.Name, "Updated Template")
	}
	if updated.Version != 2 {
		t.Errorf("Updated template Version = %d, want 2", updated.Version)
	}
}

// TestTemplateManagerUpdateVersionIncrement 测试版本自动递增
func TestTemplateManagerUpdateVersionIncrement(t *testing.T) {
	mgr := template.NewTemplateManager()

	// 创建初始版本
	tpl := createTestTemplate()
	tpl.Version = 1
	err := mgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// 更新模板(不指定版本,应该自动递增)
	updatedTpl := createTestTemplate()
	updatedTpl.Name = "Updated Template"
	// Version 应该由 Update 方法自动设置为下一个版本

	err = mgr.Update(tpl.ID, updatedTpl)
	if err != nil {
		t.Fatalf("Update() failed: %v", err)
	}

	// 验证版本已递增
	updated, err := mgr.Get(tpl.ID, 0) // 获取最新版本
	if err != nil {
		t.Fatalf("Get() failed for latest version: %v", err)
	}
	if updated.Version <= tpl.Version {
		t.Errorf("Updated template Version = %d, should be greater than %d", updated.Version, tpl.Version)
	}
}

// TestTemplateManagerUpdateNotFound 测试更新不存在的模板
func TestTemplateManagerUpdateNotFound(t *testing.T) {
	mgr := template.NewTemplateManager()

	tpl := createTestTemplate()

	// 尝试更新不存在的模板
	err := mgr.Update("non-existent", tpl)
	if err == nil {
		t.Error("Update() should fail when template does not exist")
	}
}

// TestTemplateManagerUpdateInvalidTemplate 测试更新无效模板
func TestTemplateManagerUpdateInvalidTemplate(t *testing.T) {
	mgr := template.NewTemplateManager()

	// 先创建模板
	tpl := createTestTemplate()
	err := mgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// 尝试更新为无效模板(缺少开始节点)
	invalidTpl := &template.Template{
		ID:    tpl.ID,
		Name:  "Invalid Template",
		Nodes: map[string]*template.Node{},
	}

	err = mgr.Update(tpl.ID, invalidTpl)
	if err == nil {
		t.Error("Update() should fail when template is invalid")
	}
}

// TestTemplateManagerUpdateIsolation 测试更新返回的模板是隔离的
func TestTemplateManagerUpdateIsolation(t *testing.T) {
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
	updatedTpl.Version = 2

	err = mgr.Update(tpl.ID, updatedTpl)
	if err != nil {
		t.Fatalf("Update() failed: %v", err)
	}

	// 获取更新后的模板
	got, err := mgr.Get(tpl.ID, 2)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	// 修改返回的模板
	got.Name = "Modified Name"

	// 再次获取,验证原模板未被修改
	got2, err := mgr.Get(tpl.ID, 2)
	if err != nil {
		t.Fatalf("Get() failed on second call: %v", err)
	}
	if got2.Name == "Modified Name" {
		t.Error("Get() should return a copy, modifications should not affect stored template")
	}
	if got2.Name != "Updated Template" {
		t.Errorf("Get() returned template Name = %q, want %q", got2.Name, "Updated Template")
	}
}

// TestTemplateManagerUpdateTimeStamps 测试更新时间戳
func TestTemplateManagerUpdateTimeStamps(t *testing.T) {
	mgr := template.NewTemplateManager()

	// 创建模板
	tpl := createTestTemplate()
	originalTime := time.Now().Add(-1 * time.Hour)
	tpl.CreatedAt = originalTime
	tpl.UpdatedAt = originalTime

	err := mgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// 等待一小段时间
	time.Sleep(10 * time.Millisecond)

	// 更新模板
	updatedTpl := createTestTemplate()
	updatedTpl.Name = "Updated Template"
	updatedTpl.Version = 2

	err = mgr.Update(tpl.ID, updatedTpl)
	if err != nil {
		t.Fatalf("Update() failed: %v", err)
	}

	// 验证更新时间已更新
	updated, err := mgr.Get(tpl.ID, 2)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}
	if updated.UpdatedAt.Before(originalTime) || updated.UpdatedAt.Equal(originalTime) {
		t.Error("Updated template UpdatedAt should be after original time")
	}
}

