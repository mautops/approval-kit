package template

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/mautops/approval-kit/internal/errors"
)

// memoryTemplateManager 内存实现的模板管理器
type memoryTemplateManager struct {
	mu        sync.RWMutex
	templates map[string]map[int]*Template // templateID -> version -> Template
}

// NewTemplateManager 创建新的模板管理器实例(内存实现)
func NewTemplateManager() TemplateManager {
	return &memoryTemplateManager{
		templates: make(map[string]map[int]*Template),
	}
}

// Create 创建新的审批模板
func (m *memoryTemplateManager) Create(tpl *Template) error {
	if tpl == nil {
		return fmt.Errorf("template cannot be nil")
	}

	// 验证模板
	if err := tpl.Validate(); err != nil {
		return errors.ErrInvalidTemplate
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// 检查是否已存在
	if versions, exists := m.templates[tpl.ID]; exists {
		if _, versionExists := versions[tpl.Version]; versionExists {
			return fmt.Errorf("template %q version %d already exists", tpl.ID, tpl.Version)
		}
	}

	// 创建模板副本(避免外部修改)
	templateCopy := tpl.Clone()

	// 存储模板
	if m.templates[tpl.ID] == nil {
		m.templates[tpl.ID] = make(map[int]*Template)
	}
	m.templates[tpl.ID][tpl.Version] = templateCopy

	return nil
}

// Update 更新现有审批模板
// 更新会创建新版本,原版本保持不变
func (m *memoryTemplateManager) Update(id string, tpl *Template) error {
	if tpl == nil {
		return fmt.Errorf("template cannot be nil")
	}

	// 验证模板 ID 匹配
	if tpl.ID != "" && tpl.ID != id {
		return fmt.Errorf("template ID mismatch: got %q, want %q", tpl.ID, id)
	}

	// 设置模板 ID
	tpl.ID = id

	// 验证模板
	if err := tpl.Validate(); err != nil {
		return errors.ErrInvalidTemplate
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// 检查模板是否存在
	versions, exists := m.templates[id]
	if !exists {
		return fmt.Errorf("template %q not found", id)
	}

	// 自动递增版本号(忽略传入的版本号)
	maxVersion := 0
	for v := range versions {
		if v > maxVersion {
			maxVersion = v
		}
	}
	newVersion := maxVersion + 1

	// 创建模板副本
	templateCopy := tpl.Clone()
	templateCopy.ID = id
	templateCopy.Version = newVersion
	templateCopy.UpdatedAt = time.Now()

	// 如果 CreatedAt 为零值,设置为当前时间
	if templateCopy.CreatedAt.IsZero() {
		templateCopy.CreatedAt = time.Now()
	}

	// 存储新版本
	m.templates[id][newVersion] = templateCopy

	return nil
}

// Get 获取指定版本的审批模板
func (m *memoryTemplateManager) Get(id string, version int) (*Template, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	versions, exists := m.templates[id]
	if !exists {
		return nil, fmt.Errorf("template %q not found", id)
	}

	// 如果 version 为 0,返回最新版本
	if version == 0 {
		var latestVersion int
		var latestTemplate *Template
		for v, t := range versions {
			if v > latestVersion {
				latestVersion = v
				latestTemplate = t
			}
		}
		if latestTemplate == nil {
			return nil, fmt.Errorf("template %q has no versions", id)
		}
		return latestTemplate.Clone(), nil
	}

	// 返回指定版本
	tpl, exists := versions[version]
	if !exists {
		return nil, fmt.Errorf("template %q version %d not found", id, version)
	}

	return tpl.Clone(), nil
}

// Delete 删除审批模板
// 删除会移除所有版本的模板
func (m *memoryTemplateManager) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 检查模板是否存在
	_, exists := m.templates[id]
	if !exists {
		return fmt.Errorf("template %q not found", id)
	}

	// 删除所有版本
	delete(m.templates, id)

	return nil
}

// ListVersions 列出模板的所有版本号
// 返回版本号列表,按版本号升序排列
func (m *memoryTemplateManager) ListVersions(id string) ([]int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 检查模板是否存在
	versions, exists := m.templates[id]
	if !exists {
		return nil, fmt.Errorf("template %q not found", id)
	}

	// 收集所有版本号
	versionList := make([]int, 0, len(versions))
	for v := range versions {
		versionList = append(versionList, v)
	}

	// 按版本号升序排序(使用标准库排序)
	sort.Ints(versionList)

	return versionList, nil
}
