---
layout: center
class: text-center
transition: fade
mdc: true
---

# 模板与任务管理架构

<div class="text-base text-gray-400 mt-2 mb-4">Template 结构、Task 结构、版本控制、并发安全</div>

<div class="grid grid-cols-2 gap-6 mt-4">
  <div class="text-left">
    <div class="text-lg font-bold mb-3 text-blue-400">模板管理架构</div>
    <div class="grid grid-cols-2 gap-3">
      <div class="bg-blue-500/10 rounded-lg p-3 border border-blue-500/20">
        <div class="text-lg mb-1">📋</div>
        <div class="font-bold text-blue-400 mb-1 text-sm">Template 结构</div>
        <div class="text-xs text-gray-400">ID、Name、Description、Version、Nodes、Edges、Config</div>
      </div>
      <div class="bg-green-500/10 rounded-lg p-3 border border-green-500/20">
        <div class="text-lg mb-1">🔗</div>
        <div class="font-bold text-green-400 mb-1 text-sm">节点类型</div>
        <div class="text-xs text-gray-400">Start、Approval、Condition、End</div>
      </div>
      <div class="bg-purple-500/10 rounded-lg p-3 border border-purple-500/20">
        <div class="text-lg mb-1">🔀</div>
        <div class="font-bold text-purple-400 mb-1 text-sm">节点连接</div>
        <div class="text-xs text-gray-400">Edge 定义节点间的连接关系</div>
      </div>
      <div class="bg-yellow-500/10 rounded-lg p-3 border border-yellow-500/20">
        <div class="text-lg mb-1">📌</div>
        <div class="font-bold text-yellow-400 mb-1 text-sm">版本控制</div>
        <div class="text-xs text-gray-400">模板版本管理机制</div>
      </div>
    </div>
  </div>
  
  <div class="text-left">
    <div class="text-lg font-bold mb-3 text-cyan-400">任务管理架构</div>
    <div class="grid grid-cols-2 gap-3">
      <div class="bg-cyan-500/10 rounded-lg p-3 border border-cyan-500/20">
        <div class="text-lg mb-1">✅</div>
        <div class="font-bold text-cyan-400 mb-1 text-sm">Task 结构</div>
        <div class="text-xs text-gray-400">基本信息、状态信息、运行时数据、审批记录、状态变更历史</div>
      </div>
      <div class="bg-pink-500/10 rounded-lg p-3 border border-pink-500/20">
        <div class="text-lg mb-1">🔒</div>
        <div class="font-bold text-pink-400 mb-1 text-sm">并发安全</div>
        <div class="text-xs text-gray-400">读写锁保护</div>
      </div>
      <div class="bg-orange-500/10 rounded-lg p-3 border border-orange-500/20">
        <div class="text-lg mb-1">🔄</div>
        <div class="font-bold text-orange-400 mb-1 text-sm">数据流转</div>
        <div class="text-xs text-gray-400">节点输出数据、审批人列表、审批结果</div>
      </div>
      <div class="bg-indigo-500/10 rounded-lg p-3 border border-indigo-500/20">
        <div class="text-lg mb-1">📊</div>
        <div class="font-bold text-indigo-400 mb-1 text-sm">状态管理</div>
        <div class="text-xs text-gray-400">状态信息、状态变更历史</div>
      </div>
    </div>
  </div>
</div>
