---
layout: center
class: text-center
transition: fade
mdc: true
---

# 整体架构与状态机

<div class="text-base text-gray-400 mt-2 mb-4">分层设计、状态机引擎、并发安全</div>

<div class="grid grid-cols-2 gap-6 mt-4">
  <div class="text-left">
    <div class="text-lg font-bold mb-3 text-blue-400">整体架构</div>
    <div class="bg-blue-500/10 rounded-lg p-4 border border-blue-500/20">
      <div class="text-sm font-semibold mb-2 text-blue-300">业务系统层</div>
      <div class="text-xs text-gray-400 mb-3">数据持久化、用户认证、API 接口</div>
      <div class="text-center text-gray-500 mb-3">↓</div>
      <div class="text-sm font-semibold mb-2 text-green-300">审批流核心库</div>
      <div class="text-xs text-gray-400 mb-2">• 状态机引擎</div>
      <div class="text-xs text-gray-400 mb-2">• 模板管理</div>
      <div class="text-xs text-gray-400 mb-2">• 任务管理</div>
      <div class="text-xs text-gray-400 mb-2">• 节点执行引擎</div>
      <div class="text-xs text-gray-400 mb-3">• 事件通知</div>
      <div class="text-center text-gray-500 mb-3">↓</div>
      <div class="text-sm font-semibold mb-2 text-purple-300">扩展接口层</div>
      <div class="text-xs text-gray-400">HTTP 客户端、条件评估、存储接口</div>
    </div>
  </div>
  
  <div class="text-left">
    <div class="text-lg font-bold mb-3 text-cyan-400">状态机设计</div>
    <div class="grid grid-cols-2 gap-3">
      <div class="bg-cyan-500/10 rounded-lg p-3 border border-cyan-500/20">
        <div class="text-lg mb-1">📊</div>
        <div class="font-bold text-cyan-400 mb-1 text-sm">状态定义</div>
        <div class="text-xs text-gray-400">pending、submitted、approving、approved、rejected、cancelled、timeout</div>
      </div>
      <div class="bg-pink-500/10 rounded-lg p-3 border border-pink-500/20">
        <div class="text-lg mb-1">🔄</div>
        <div class="font-bold text-pink-400 mb-1 text-sm">状态转换</div>
        <div class="text-xs text-gray-400">CanTransition、Transition、GetValidTransitions</div>
      </div>
      <div class="bg-orange-500/10 rounded-lg p-3 border border-orange-500/20">
        <div class="text-lg mb-1">🔒</div>
        <div class="font-bold text-orange-400 mb-1 text-sm">并发安全</div>
        <div class="text-xs text-gray-400">版本号机制、原子操作</div>
      </div>
      <div class="bg-indigo-500/10 rounded-lg p-3 border border-indigo-500/20">
        <div class="text-lg mb-1">✅</div>
        <div class="font-bold text-indigo-400 mb-1 text-sm">转换规则</div>
        <div class="text-xs text-gray-400">确保状态转换的合法性</div>
      </div>
    </div>
  </div>
</div>
