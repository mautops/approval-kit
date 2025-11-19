package template_test

import (
	"testing"

	"github.com/mautops/approval-kit/internal/template"
)

// TestTemplateManagerListVersions 测试模板版本列表功能
func TestTemplateManagerListVersions(t *testing.T) {
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

	// 获取版本列表
	versions, err := mgr.ListVersions(tpl1.ID)
	if err != nil {
		t.Fatalf("ListVersions() failed: %v", err)
	}

	// 验证版本数量
	if len(versions) != 3 {
		t.Errorf("ListVersions() returned %d versions, want 3", len(versions))
	}

	// 验证版本号包含所有版本
	versionMap := make(map[int]bool)
	for _, v := range versions {
		versionMap[v] = true
	}
	if !versionMap[1] || !versionMap[2] || !versionMap[3] {
		t.Error("ListVersions() should return all versions")
	}
}

// TestTemplateManagerListVersionsOrdered 测试版本列表排序
func TestTemplateManagerListVersionsOrdered(t *testing.T) {
	mgr := template.NewTemplateManager()

	// 创建多个版本(不按顺序)
	tpl3 := createTestTemplate()
	tpl3.Version = 3
	err := mgr.Create(tpl3)
	if err != nil {
		t.Fatalf("Create() failed for version 3: %v", err)
	}

	tpl1 := createTestTemplate()
	tpl1.Version = 1
	err = mgr.Create(tpl1)
	if err != nil {
		t.Fatalf("Create() failed for version 1: %v", err)
	}

	tpl2 := createTestTemplate()
	tpl2.Version = 2
	err = mgr.Create(tpl2)
	if err != nil {
		t.Fatalf("Create() failed for version 2: %v", err)
	}

	// 获取版本列表
	versions, err := mgr.ListVersions(tpl1.ID)
	if err != nil {
		t.Fatalf("ListVersions() failed: %v", err)
	}

	// 验证版本列表按升序排列
	if len(versions) < 3 {
		t.Fatalf("ListVersions() returned %d versions, want at least 3", len(versions))
	}

	for i := 1; i < len(versions); i++ {
		if versions[i-1] > versions[i] {
			t.Errorf("ListVersions() returned unsorted versions: %v", versions)
			break
		}
	}
}

// TestTemplateManagerListVersionsNotFound 测试查询不存在的模板版本列表
func TestTemplateManagerListVersionsNotFound(t *testing.T) {
	mgr := template.NewTemplateManager()

	// 查询不存在的模板
	_, err := mgr.ListVersions("non-existent")
	if err == nil {
		t.Error("ListVersions() should fail when template does not exist")
	}
}

// TestTemplateManagerListVersionsEmpty 测试空版本列表
func TestTemplateManagerListVersionsEmpty(t *testing.T) {
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

	// 查询已删除的模板版本列表
	versions, err := mgr.ListVersions(tpl.ID)
	if err == nil {
		// 如果返回成功,版本列表应该为空
		if len(versions) != 0 {
			t.Errorf("ListVersions() returned %d versions for deleted template, want 0", len(versions))
		}
	}
}

// TestTemplateManagerListVersionsAfterUpdate 测试更新后的版本列表
func TestTemplateManagerListVersionsAfterUpdate(t *testing.T) {
	mgr := template.NewTemplateManager()

	// 创建模板
	tpl := createTestTemplate()
	err := mgr.Create(tpl)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// 更新模板(创建新版本)
	updatedTpl := createTestTemplate()
	updatedTpl.Name = "Updated Template"
	err = mgr.Update(tpl.ID, updatedTpl)
	if err != nil {
		t.Fatalf("Update() failed: %v", err)
	}

	// 获取版本列表
	versions, err := mgr.ListVersions(tpl.ID)
	if err != nil {
		t.Fatalf("ListVersions() failed: %v", err)
	}

	// 验证有两个版本
	if len(versions) != 2 {
		t.Errorf("ListVersions() returned %d versions, want 2", len(versions))
	}

	// 验证版本号
	versionMap := make(map[int]bool)
	for _, v := range versions {
		versionMap[v] = true
	}
	if !versionMap[1] || !versionMap[2] {
		t.Error("ListVersions() should return versions 1 and 2")
	}
}

