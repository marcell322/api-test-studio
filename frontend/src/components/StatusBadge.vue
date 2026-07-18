<template>
  <span class="px-2 py-0.5 rounded text-xs font-mono font-semibold" :class="colorClass">{{ display }}</span>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({ status: { type: Number, required: true } })

const display = computed(() => (props.status === 0 ? 'ERR' : props.status))

const colorClass = computed(() => {
  const s = props.status
  if (s === 0) return 'bg-servererr/20 text-servererr'
  if (s < 300) return 'bg-ok/20 text-ok'
  if (s < 400) return 'bg-redirect/20 text-redirect'
  if (s < 500) return 'bg-clienterr/20 text-clienterr'
  return 'bg-servererr/20 text-servererr'
})
</script>
