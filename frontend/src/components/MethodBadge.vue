<template>
  <span class="px-2 py-0.5 rounded text-xs font-mono font-semibold" :class="colorClass">{{ method }}</span>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({ method: { type: String, required: true } })

const colorClass = computed(() => {
  switch (props.method?.toUpperCase()) {
    case 'GET':
      return 'bg-get/20 text-get'
    case 'POST':
      return 'bg-post/20 text-post'
    case 'PUT':
      return 'bg-put/20 text-put'
    case 'PATCH':
      return 'bg-patch/20 text-patch'
    case 'DELETE':
      return 'bg-del/20 text-del'
    default:
      return 'bg-muted/20 text-muted'
  }
})
</script>
