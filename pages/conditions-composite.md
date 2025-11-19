---
layout: center
class: text-center
transition: fade
mdc: true
---

# 条件节点 - 组合条件与数据传递

<div class="text-base text-gray-400 mt-2 mb-4">AND、OR 逻辑组合、节点输出、数据读取、上下文缓存</div>

<div class="mt-4 max-w-5xl mx-auto">
  <div class="bg-gradient-to-r from-blue-500/20 to-purple-500/20 rounded-xl p-5 border border-blue-500/30 mb-4">
    <div class="flex items-center justify-center mb-3">
      <div class="w-10 h-10 rounded-full bg-blue-500 flex items-center justify-center text-white font-bold mr-3">1</div>
      <div class="text-lg font-bold text-blue-400">组合条件</div>
    </div>
    <div class="flex justify-center gap-3 flex-wrap">
      <div class="px-4 py-2 bg-blue-500/30 rounded-lg border border-blue-500/50">
        <span class="text-sm font-semibold text-blue-200">AND</span>
        <span class="text-xs text-gray-300 ml-2">逻辑与</span>
      </div>
      <div class="px-4 py-2 bg-purple-500/30 rounded-lg border border-purple-500/50">
        <span class="text-sm font-semibold text-purple-200">OR</span>
        <span class="text-xs text-gray-300 ml-2">逻辑或</span>
      </div>
    </div>
    <div class="text-xs text-gray-400 mt-3">支持复杂的业务规则组合</div>
  </div>

  <div class="flex items-center justify-center mb-4">
    <div class="w-0.5 h-8 bg-gradient-to-b from-blue-500 to-green-500"></div>
  </div>

  <div class="grid grid-cols-3 gap-3">
    <div class="bg-green-500/10 rounded-lg p-3 border border-green-500/20">
      <div class="text-lg mb-1">📤</div>
      <div class="font-bold text-green-400 mb-1 text-sm">节点输出</div>
      <div class="text-xs text-gray-400">每个节点可以输出 JSON 格式的数据</div>
    </div>
    <div class="bg-cyan-500/10 rounded-lg p-3 border border-cyan-500/20">
      <div class="text-lg mb-1">📥</div>
      <div class="font-bold text-cyan-400 mb-1 text-sm">数据读取</div>
      <div class="text-xs text-gray-400">后续节点可以读取前面节点的输出</div>
    </div>
    <div class="bg-orange-500/10 rounded-lg p-3 border border-orange-500/20">
      <div class="text-lg mb-1">⚡</div>
      <div class="font-bold text-orange-400 mb-1 text-sm">上下文缓存</div>
      <div class="text-xs text-gray-400">避免重复查询,提高性能</div>
    </div>
  </div>
</div>
