---
layout: center
class: text-center
transition: fade
mdc: true
---

# 事件通知

<div class="text-base text-gray-400 mt-2 mb-4">Webhook 异步推送、重试机制、幂等性保证</div>

<div class="mt-4 max-w-5xl mx-auto">
  <div class="bg-indigo-500/10 rounded-xl p-4 border border-indigo-500/20 mb-3">
    <div class="text-base font-bold text-indigo-400 mb-2">事件类型</div>
    <div class="flex flex-wrap gap-2 justify-center">
      <span class="px-2.5 py-1 bg-indigo-500/20 rounded text-xs text-indigo-300">任务创建</span>
      <span class="px-2.5 py-1 bg-indigo-500/20 rounded text-xs text-indigo-300">任务提交</span>
      <span class="px-2.5 py-1 bg-indigo-500/20 rounded text-xs text-indigo-300">节点激活</span>
      <span class="px-2.5 py-1 bg-indigo-500/20 rounded text-xs text-indigo-300">审批操作</span>
      <span class="px-2.5 py-1 bg-indigo-500/20 rounded text-xs text-indigo-300">任务完成</span>
    </div>
  </div>
  
  <div class="grid grid-cols-4 gap-2">
    <div class="bg-purple-500/10 rounded-lg p-2.5 border border-purple-500/20">
      <div class="text-sm mb-1">🔗</div>
      <div class="font-bold text-purple-400 mb-0.5 text-xs">Webhook</div>
      <div class="text-xs text-gray-400">支持多个地址、认证配置</div>
    </div>
    <div class="bg-cyan-500/10 rounded-lg p-2.5 border border-cyan-500/20">
      <div class="text-sm mb-1">⚡</div>
      <div class="font-bold text-cyan-400 mb-0.5 text-xs">异步推送</div>
      <div class="text-xs text-gray-400">不阻塞主流程</div>
    </div>
    <div class="bg-orange-500/10 rounded-lg p-2.5 border border-orange-500/20">
      <div class="text-sm mb-1">🔄</div>
      <div class="font-bold text-orange-400 mb-0.5 text-xs">重试机制</div>
      <div class="text-xs text-gray-400">指数退避策略</div>
    </div>
    <div class="bg-green-500/10 rounded-lg p-2.5 border border-green-500/20">
      <div class="text-sm mb-1">🔒</div>
      <div class="font-bold text-green-400 mb-0.5 text-xs">幂等性</div>
      <div class="text-xs text-gray-400">确保不重复处理</div>
    </div>
  </div>
</div>
