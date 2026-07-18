<template>
  <div class="flex-1 overflow-y-auto p-6">
    <h2 class="font-display text-lg font-semibold mb-4">History</h2>
    <div v-if="loading" class="text-sm text-muted font-mono">Loading…</div>
    <div v-else-if="!items.length" class="text-sm text-muted">No requests sent yet. Send one to see it here.</div>
    <div class="space-y-1">
      <div
        v-for="item in items"
        :key="item.id"
        class="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-surface2 cursor-pointer group"
        @click="emit('replay', item)"
      >
        <MethodBadge :method="item.method" />
        <StatusBadge :status="item.status_code" />
        <span class="flex-1 truncate text-sm font-mono text-muted">{{ item.url }}</span>
        <span class="text-xs text-muted font-mono">{{ item.response_time_ms }}ms</span>
        <span class="text-xs text-muted">{{ formatTime(item.created_at) }}</span>
        <button @click.stop="remove(item)" class="hidden group-hover:block text-xs text-muted hover:text-del">✕</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import client from '../api/client'
import MethodBadge from './MethodBadge.vue'
import StatusBadge from './StatusBadge.vue'

const emit = defineEmits(['replay'])
const items = ref([])
const loading = ref(true)

function formatTime(ts) {
  return new Date(ts).toLocaleString()
}

async function load() {
  loading.value = true
  try {
    const res = await client.get('/history')
    items.value = res.data.data || []
  } finally {
    loading.value = false
  }
}

async function remove(item) {
  await client.delete(`/history/${item.id}`)
  items.value = items.value.filter((x) => x.id !== item.id)
}

onMounted(load)
defineExpose({ load })
</script>
