---
layout: center
class: text-center
transition: fade
mdc: true
---

# 节点执行引擎与事件通知

<div class="text-base text-gray-400 mt-2 mb-4">NodeExecutor 接口、节点上下文、异步推送、重试机制</div>

<div class="grid grid-cols-2 gap-6 mt-4">
  <div class="text-left">
    <div class="text-lg font-bold mb-3 text-blue-400">节点执行引擎</div>
    <div class="grid grid-cols-2 gap-3">
      <div class="bg-blue-500/10 rounded-lg p-3 border border-blue-500/20">
        <div class="text-lg mb-1">🔌</div>
        <div class="font-bold text-blue-400 mb-1 text-sm">NodeExecutor 接口</div>
        <div class="text-xs text-gray-400">统一的节点执行接口</div>
      </div>
      <div class="bg-green-500/10 rounded-lg p-3 border border-green-500/20">
        <div class="text-lg mb-1">⚙️</div>
        <div class="font-bold text-green-400 mb-1 text-sm">节点类型执行器</div>
        <div class="text-xs text-gray-400">StartNodeExecutor、ApprovalNodeExecutor、ConditionNodeExecutor、EndNodeExecutor</div>
      </div>
      <div class="bg-purple-500/10 rounded-lg p-3 border border-purple-500/20">
        <div class="text-lg mb-1">📋</div>
        <div class="font-bold text-purple-400 mb-1 text-sm">节点上下文</div>
        <div class="text-xs text-gray-400">NodeContext 提供执行上下文信息</div>
      </div>
      <div class="bg-yellow-500/10 rounded-lg p-3 border border-yellow-500/20">
        <div class="text-lg mb-1">🔄</div>
        <div class="font-bold text-yellow-400 mb-1 text-sm">数据传递</div>
        <div class="text-xs text-gray-400">节点输出数据、上下文缓存</div>
      </div>
    </div>
  </div>
  
  <div class="text-left">
    <div class="text-lg font-bold mb-3 text-cyan-400">事件通知机制</div>
    <div class="grid grid-cols-2 gap-3">
      <div class="bg-cyan-500/10 rounded-lg p-3 border border-cyan-500/20">
        <div class="text-lg mb-1">📢</div>
        <div class="font-bold text-cyan-400 mb-1 text-sm">事件类型</div>
        <div class="text-xs text-gray-400">任务创建、提交、节点激活、审批操作、任务完成等</div>
      </div>
      <div class="bg-pink-500/10 rounded-lg p-3 border border-pink-500/20">
        <div class="text-lg mb-1">⚡</div>
        <div class="font-bold text-pink-400 mb-1 text-sm">异步推送</div>
        <div class="text-xs text-gray-400">goroutine + channel 实现</div>
      </div>
      <div class="bg-orange-500/10 rounded-lg p-3 border border-orange-500/20">
        <div class="text-lg mb-1">🔄</div>
        <div class="font-bold text-orange-400 mb-1 text-sm">重试机制</div>
        <div class="text-xs text-gray-400">指数退避策略</div>
      </div>
      <div class="bg-indigo-500/10 rounded-lg p-3 border border-indigo-500/20">
        <div class="text-lg mb-1">🔗</div>
        <div class="font-bold text-indigo-400 mb-1 text-sm">Webhook 配置</div>
        <div class="text-xs text-gray-400">支持多个 Webhook 地址</div>
      </div>
    </div>
  </div>
</div>
