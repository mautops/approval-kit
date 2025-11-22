package template

import (
	internalTemplate "github.com/mautops/approval-kit/internal/template"
)

// Template 表示审批模板
// 模板定义完整的审批流程结构,包括节点定义和连接关系
// 与 internal/template.Template 结构相同,但位于 pkg 目录,可以被外部导入
type Template = internalTemplate.Template

// TemplateConfig 模板全局配置
// 与 internal/template.TemplateConfig 结构相同,但位于 pkg 目录,可以被外部导入
type TemplateConfig = internalTemplate.TemplateConfig

// WebhookConfig Webhook 配置
// 与 internal/template.WebhookConfig 结构相同,但位于 pkg 目录,可以被外部导入
type WebhookConfig = internalTemplate.WebhookConfig

// AuthConfig 认证配置
// 与 internal/template.AuthConfig 结构相同,但位于 pkg 目录,可以被外部导入
type AuthConfig = internalTemplate.AuthConfig

// TemplateFromInternal 将 internal.Template 转换为 pkg.Template
func TemplateFromInternal(t *internalTemplate.Template) *Template {
	return (*Template)(t)
}

// TemplateToInternal 将 pkg.Template 转换为 internal.Template
func TemplateToInternal(t *Template) *internalTemplate.Template {
	return (*internalTemplate.Template)(t)
}

