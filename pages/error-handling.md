---
layout: center
class: text-center
transition: fade
mdc: true
---

# 错误处理

<div class="text-base text-gray-400 mt-2 mb-4">错误类型、错误处理、调试技巧</div>

<div class="mt-4 max-w-4xl mx-auto">
  <div class="bg-red-500/10 rounded-xl p-4 border border-red-500/20 mb-3">
    <div class="text-base font-bold text-red-400 mb-2">错误类型</div>
    <div class="flex flex-wrap gap-2 justify-center">
      <span class="px-3 py-1.5 bg-red-500/20 rounded-lg text-xs text-red-300">无效模板</span>
      <span class="px-3 py-1.5 bg-red-500/20 rounded-lg text-xs text-red-300">状态转换错误</span>
      <span class="px-3 py-1.5 bg-red-500/20 rounded-lg text-xs text-red-300">节点未找到</span>
      <span class="px-3 py-1.5 bg-red-500/20 rounded-lg text-xs text-red-300">审批人未找到</span>
      <span class="px-3 py-1.5 bg-red-500/20 rounded-lg text-xs text-red-300">并发修改</span>
    </div>
  </div>
  
  <div class="grid grid-cols-2 gap-3">
    <div class="bg-orange-500/10 rounded-lg p-3 border border-orange-500/20">
      <div class="text-sm font-bold text-orange-400 mb-1">错误处理</div>
      <div class="text-xs text-gray-400 text-left space-y-1">
        <div>• 使用 error 接口返回错误</div>
        <div>• 错误包装和错误链追踪</div>
        <div>• 错误分类和处理</div>
      </div>
    </div>
    <div class="bg-yellow-500/10 rounded-lg p-3 border border-yellow-500/20">
      <div class="text-sm font-bold text-yellow-400 mb-1">调试技巧</div>
      <div class="text-xs text-gray-400 text-left space-y-1">
        <div>• 查看状态变更历史</div>
        <div>• 检查审批记录</div>
        <div>• 问题排查方法</div>
      </div>
    </div>
  </div>
</div>
