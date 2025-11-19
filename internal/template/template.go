package template

import (
	"time"
)

// Template 表示审批模板
// 模板定义完整的审批流程结构,包括节点定义和连接关系
type Template struct {
	// 基本信息
	ID          string    // 模板 ID
	Name        string    // 模板名称
	Description string    // 模板描述
	Version     int       // 版本号
	CreatedAt   time.Time // 创建时间
	UpdatedAt   time.Time // 更新时间

	// 节点定义
	// 键为节点 ID,值为节点定义
	Nodes map[string]*Node

	// 节点连接关系 (有向图)
	// 定义节点间的连接,支持条件分支
	Edges []*Edge

	// 全局配置
	// 包含 Webhook 配置、超时配置等
	Config *TemplateConfig
}

// TemplateConfig 模板全局配置
type TemplateConfig struct {
	// Webhook 配置
	Webhooks []*WebhookConfig

	// 其他全局配置可以在这里扩展
}

// WebhookConfig Webhook 配置
type WebhookConfig struct {
	URL     string            // Webhook 地址
	Method  string            // 请求方法 (GET/POST)
	Headers map[string]string // 请求头
	Auth    *AuthConfig       // 认证配置
}

// AuthConfig 认证配置
type AuthConfig struct {
	Type  string // 认证类型 (token/signature)
	Token string // Token 值
	Key   string // 签名密钥
}
