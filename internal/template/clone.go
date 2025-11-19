package template

// Clone 创建模板的深拷贝
func (t *Template) Clone() *Template {
	if t == nil {
		return nil
	}

	clone := &Template{
		ID:          t.ID,
		Name:        t.Name,
		Description: t.Description,
		Version:     t.Version,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}

	// 复制 Nodes
	if t.Nodes != nil {
		clone.Nodes = make(map[string]*Node, len(t.Nodes))
		for k, v := range t.Nodes {
			clone.Nodes[k] = v.Clone()
		}
	}

	// 复制 Edges
	if t.Edges != nil {
		clone.Edges = make([]*Edge, len(t.Edges))
		for i, e := range t.Edges {
			clone.Edges[i] = e.Clone()
		}
	}

	// 复制 Config
	if t.Config != nil {
		clone.Config = t.Config.Clone()
	}

	return clone
}

// Clone 创建节点的深拷贝
func (n *Node) Clone() *Node {
	if n == nil {
		return nil
	}

	clone := &Node{
		ID:    n.ID,
		Name:  n.Name,
		Type:  n.Type,
		Order: n.Order,
	}

	// 注意: Config 是接口,这里只复制引用
	// 如果需要深拷贝 Config,需要根据具体类型实现
	clone.Config = n.Config

	return clone
}

// Clone 创建边的深拷贝
func (e *Edge) Clone() *Edge {
	if e == nil {
		return nil
	}

	return &Edge{
		From:      e.From,
		To:        e.To,
		Condition: e.Condition,
	}
}

// Clone 创建模板配置的深拷贝
func (c *TemplateConfig) Clone() *TemplateConfig {
	if c == nil {
		return nil
	}

	clone := &TemplateConfig{}

	// 复制 Webhooks
	if c.Webhooks != nil {
		clone.Webhooks = make([]*WebhookConfig, len(c.Webhooks))
		for i, w := range c.Webhooks {
			clone.Webhooks[i] = w.Clone()
		}
	}

	return clone
}

// Clone 创建 Webhook 配置的深拷贝
func (w *WebhookConfig) Clone() *WebhookConfig {
	if w == nil {
		return nil
	}

	clone := &WebhookConfig{
		URL:    w.URL,
		Method: w.Method,
	}

	// 复制 Headers
	if w.Headers != nil {
		clone.Headers = make(map[string]string, len(w.Headers))
		for k, v := range w.Headers {
			clone.Headers[k] = v
		}
	}

	// 复制 Auth
	if w.Auth != nil {
		clone.Auth = &AuthConfig{
			Type:  w.Auth.Type,
			Token: w.Auth.Token,
			Key:   w.Auth.Key,
		}
	}

	return clone
}

