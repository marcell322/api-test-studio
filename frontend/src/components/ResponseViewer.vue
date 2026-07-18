<template>
  <div class="border-t border-border p-4 max-h-64 overflow-y-auto">
    <div class="flex items-center gap-3 mb-3">
      <StatusBadge :status="response.status_code" />
      <span class="text-xs text-muted font-mono">{{ response.response_time_ms }}ms</span>
    </div>
    <div v-if="response.headers && Object.keys(response.headers).length" class="mb-3">
      <div class="text-xs text-muted uppercase tracking-wide mb-1 font-mono">Response Headers</div>
      <div class="font-mono text-xs text-muted space-y-0.5">
        <div v-for="(v, k) in response.headers" :key="k"><span class="text-text">{{ k }}</span>: {{ v }}</div>
      </div>
    </div>
    <div>
      <div class="text-xs text-muted uppercase tracking-wide mb-1 font-mono">Body</div>
      <pre class="bg-surface border border-border rounded-lg p-3 text-xs font-mono overflow-x-auto whitespace-pre-wrap break-all">{{ prettyBody }}</pre>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import StatusBadge from './StatusBadge.vue'

const props = defineProps({ response: { type: Object, required: true } })

const prettyBody = computed(() => {
  try {
    return JSON.stringify(JSON.parse(props.response.body), null, 2)
  } catch {
    return props.response.body
  }
})
</script>
