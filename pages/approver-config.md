---
layout: center
class: text-center
transition: fade
mdc: true
---

# 审批人配置

<div class="text-base text-gray-400 mt-2 mb-4">固定审批人、动态审批人、HTTP API 配置</div>

<div class="grid grid-cols-2 gap-6 mt-4">
  <div class="text-left">
    <div class="flex items-center mb-3">
      <div class="w-1 h-8 bg-blue-500 rounded-full mr-3"></div>
      <div class="text-lg font-bold text-blue-400">固定审批人</div>
    </div>
    <div class="bg-blue-500/10 rounded-lg p-4 border border-blue-500/20">
      <div class="text-sm font-semibold text-blue-300 mb-2">FixedApproverConfig</div>
      <div class="text-xs text-gray-400 leading-relaxed">适用于审批人固定的场景,在模板中预设审批人列表,审批人必须是用户 ID,必须具体到人,不能是角色、部门等标识</div>
    </div>
  </div>
  
  <div class="text-left">
    <div class="flex items-center mb-3">
      <div class="w-1 h-8 bg-green-500 rounded-full mr-3"></div>
      <div class="text-lg font-bold text-green-400">动态审批人</div>
    </div>
    <div class="bg-green-500/10 rounded-lg p-4 border border-green-500/20">
      <div class="text-sm font-semibold text-green-300 mb-3">DynamicApproverConfig</div>
      <div class="space-y-2">
        <div class="flex items-start">
          <div class="text-xs font-semibold text-green-400 mr-2 min-w-[80px]">HTTP API:</div>
          <div class="text-xs text-gray-400">URL、Method、Headers、参数映射、响应解析</div>
        </div>
        <div class="flex items-start">
          <div class="text-xs font-semibold text-green-400 mr-2 min-w-[80px]">获取时机:</div>
          <div class="text-xs text-gray-400">任务创建时 / 节点激活时</div>
        </div>
        <div class="flex items-start">
          <div class="text-xs font-semibold text-green-400 mr-2 min-w-[80px]">重试机制:</div>
          <div class="text-xs text-gray-400">API 调用失败自动重试</div>
        </div>
      </div>
    </div>
  </div>
</div>
