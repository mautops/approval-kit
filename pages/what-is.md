---
layout: center
class: text-center
transition: fade
mdc: true
---

# 什么是 Approval Kit?

<div class="text-base text-gray-400 mt-2 mb-4">Go 语言编写的审批流核心库,专注于管理审批模板和审批任务的状态流转</div>

<div class="text-xl font-bold mt-4 mb-3">核心特性概览</div>

<div class="grid grid-cols-3 gap-4">
  <div class="bg-blue-500/10 rounded-lg p-4 border border-blue-500/20">
    <div class="text-xl mb-1">📋</div>
    <div class="font-bold text-blue-400 mb-1 text-sm">模板管理</div>
    <div class="text-xs text-gray-400">创建、更新、删除、查询、版本控制</div>
  </div>
  <div class="bg-green-500/10 rounded-lg p-4 border border-green-500/20">
    <div class="text-xl mb-1">✅</div>
    <div class="font-bold text-green-400 mb-1 text-sm">任务管理</div>
    <div class="text-xs text-gray-400">创建、提交、审批、查询</div>
  </div>
  <div class="bg-purple-500/10 rounded-lg p-4 border border-purple-500/20">
    <div class="text-xl mb-1">🔄</div>
    <div class="font-bold text-purple-400 mb-1 text-sm">状态机</div>
    <div class="text-xs text-gray-400">状态流转管理,确保转换合法性</div>
  </div>
</div>

<div class="grid grid-cols-3 gap-4 mt-3">
  <div class="bg-yellow-500/10 rounded-lg p-4 border border-yellow-500/20">
    <div class="text-xl mb-1">👥</div>
    <div class="font-bold text-yellow-400 mb-1 text-sm">多种审批模式</div>
    <div class="text-xs text-gray-400">单人、会签、或签、比例会签、顺序审批</div>
  </div>
  <div class="bg-pink-500/10 rounded-lg p-4 border border-pink-500/20">
    <div class="text-xl mb-1">🔗</div>
    <div class="font-bold text-pink-400 mb-1 text-sm">动态审批人</div>
    <div class="text-xs text-gray-400">HTTP API 动态获取审批人列表</div>
  </div>
  <div class="bg-cyan-500/10 rounded-lg p-4 border border-cyan-500/20">
    <div class="text-xl mb-1">🔀</div>
    <div class="font-bold text-cyan-400 mb-1 text-sm">条件分支</div>
    <div class="text-xs text-gray-400">数值、字符串、枚举、组合条件</div>
  </div>
</div>

<div class="grid grid-cols-3 gap-4 mt-3">
  <div class="bg-orange-500/10 rounded-lg p-4 border border-orange-500/20">
    <div class="text-xl mb-1">⚙️</div>
    <div class="font-bold text-orange-400 mb-1 text-sm">高级操作</div>
    <div class="text-xs text-gray-400">转交、加签、减签、撤回、取消</div>
  </div>
  <div class="bg-indigo-500/10 rounded-lg p-4 border border-indigo-500/20">
    <div class="text-xl mb-1">📢</div>
    <div class="font-bold text-indigo-400 mb-1 text-sm">事件通知</div>
    <div class="text-xs text-gray-400">Webhook 异步推送,支持重试机制</div>
  </div>
  <div class="bg-red-500/10 rounded-lg p-4 border border-red-500/20">
    <div class="text-xl mb-1">🔒</div>
    <div class="font-bold text-red-400 mb-1 text-sm">并发安全</div>
    <div class="text-xs text-gray-400">读写锁、版本号机制保证并发安全</div>
  </div>
</div>
