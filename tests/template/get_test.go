package template_test

import (
	"testing"

	"github.com/mautops/approval-kit/internal/template"
)

// TestTemplateManagerGet 测试模板查询功能
func TestTemplateManagerGet(t *testing.T) {
	mgr := template.NewTemplateManager()

	tpl := createTestTemplate()

	// 先创建模板
	err := mgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// 测试获取指定版本
	got, err := mgr.Get(tpl.ID, tpl.Version)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}
	if got.ID != tpl.ID {
		t.Errorf("Get() returned template ID = %q, want %q", got.ID, tpl.ID)
	}
	if got.Version != tpl.Version {
		t.Errorf("Get() returned template Version = %d, want %d", got.Version, tpl.Version)
	}
}

// TestTemplateManagerGetLatestVersion 测试获取最新版本
func TestTemplateManagerGetLatestVersion(t *testing.T) {
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

	tpl3 := createTestTemplate()
	tpl3.Version = 3
	err = mgr.Create(tpl3)
	if err != nil {
		t.Fatalf("Create() failed for version 3: %v", err)
	}

	// 测试获取最新版本(version = 0)
	got, err := mgr.Get(tpl1.ID, 0)
	if err != nil {
		t.Fatalf("Get() failed for latest version: %v", err)
	}
	if got.Version != 3 {
		t.Errorf("Get() returned template Version = %d, want 3 (latest)", got.Version)
	}
}

// TestTemplateManagerGetNotFound 测试查询不存在的模板
func TestTemplateManagerGetNotFound(t *testing.T) {
	mgr := template.NewTemplateManager()

	// 测试查询不存在的模板
	_, err := mgr.Get("non-existent", 1)
	if err == nil {
		t.Error("Get() should fail when template does not exist")
	}
}

// TestTemplateManagerGetVersionNotFound 测试查询不存在的版本
func TestTemplateManagerGetVersionNotFound(t *testing.T) {
	mgr := template.NewTemplateManager()

	tpl := createTestTemplate()
	tpl.Version = 1

	// 创建模板
	err := mgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// 测试查询不存在的版本
	_, err = mgr.Get(tpl.ID, 999)
	if err == nil {
		t.Error("Get() should fail when version does not exist")
	}
}

// TestTemplateManagerGetIsolation 测试查询返回的模板是隔离的
func TestTemplateManagerGetIsolation(t *testing.T) {
	mgr := template.NewTemplateManager()

	tpl := createTestTemplate()
	err := mgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// 获取模板
	got, err := mgr.Get(tpl.ID, tpl.Version)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	// 修改返回的模板
	got.Name = "Modified Name"

	// 再次获取,验证原模板未被修改
	got2, err := mgr.Get(tpl.ID, tpl.Version)
	if err != nil {
		t.Fatalf("Get() failed on second call: %v", err)
	}
	if got2.Name == "Modified Name" {
		t.Error("Get() should return a copy, modifications should not affect stored template")
	}
	if got2.Name != tpl.Name {
		t.Errorf("Get() returned template Name = %q, want %q", got2.Name, tpl.Name)
	}
}

