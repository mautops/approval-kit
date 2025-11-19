package template

// TemplateManager 模板管理接口
// 负责审批模板的创建、更新、查询和删除操作
// 遵循接口隔离原则,提供小而专注的接口
type TemplateManager interface {
	// Create 创建新的审批模板
	// template: 待创建的模板对象
	// 返回: 错误信息(如果模板无效或已存在)
	Create(template *Template) error

	// Update 更新现有审批模板
	// id: 模板 ID
	// template: 更新后的模板对象
	// 返回: 错误信息(如果模板不存在或无效)
	// 注意: 更新会创建新版本,原版本保持不变
	Update(id string, template *Template) error

	// Get 获取指定版本的审批模板
	// id: 模板 ID
	// version: 模板版本号,0 表示获取最新版本
	// 返回: 模板对象和错误信息
	Get(id string, version int) (*Template, error)

	// Delete 删除审批模板
	// id: 模板 ID
	// 返回: 错误信息(如果模板不存在)
	// 注意: 删除会移除所有版本的模板
	Delete(id string) error

	// ListVersions 列出模板的所有版本号
	// id: 模板 ID
	// 返回: 版本号列表(按版本号升序排列)和错误信息
	ListVersions(id string) ([]int, error)
}

