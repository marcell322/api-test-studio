<template>
  <div class="flex flex-col h-full">
    <!-- name + save row -->
    <div class="border-b border-border p-4 flex items-center gap-3">
      <input
        v-model="name"
        placeholder="Untitled request"
        class="flex-1 bg-transparent font-display text-lg font-medium focus:outline-none placeholder:text-muted"
      />
      <select v-model="collectionId" class="bg-surface border border-border rounded-lg px-2 py-1.5 text-sm focus:outline-none">
        <option :value="null">No collection</option>
        <option v-for="c in collections" :key="c.id" :value="c.id">{{ c.name }}</option>
      </select>
      <button
        @click="save"
        :disabled="!name || !collectionId || saving"
        class="text-sm border border-border rounded-lg px-3 py-1.5 hover:bg-surface2 disabled:opacity-40 disabled:cursor-not-allowed"
      >
        {{ savedRequestId ? 'Update' : 'Save' }}
      </button>
    </div>

    <!-- method + url row -->
    <div class="p-4 flex items-center gap-2 border-b border-border">
      <select
        v-model="method"
        class="bg-surface border border-border rounded-lg px-2 py-2 text-sm font-mono font-semibold focus:outline-none"
        :class="methodTextClass"
      >
        <option v-for="m in methods" :key="m" :value="m">{{ m }}</option>
      </select>
      <input
        v-model="url"
        placeholder="https://api.example.com/endpoint"
        class="flex-1 bg-surface border border-border rounded-lg px-3 py-2 text-sm font-mono focus:outline-none focus:ring-2 focus:ring-accent/50 focus:border-accent"
        @keyup.enter="send"
      />
      <button
        @click="send"
        :disabled="!url || sending"
        class="bg-accent hover:bg-accent-dim transition-colors rounded-lg px-5 py-2 text-sm font-medium text-bg disabled:opacity-50"
      >
        {{ sending ? 'Sending…' : 'Send' }}
      </button>
    </div>

    <!-- headers / body -->
    <div class="flex-1 overflow-y-auto p-4 grid grid-cols-2 gap-4 min-h-0">
      <div class="flex flex-col min-h-0">
        <div class="text-xs text-muted uppercase tracking-wide mb-2 font-mono">Headers</div>
        <div class="space-y-2">
          <div v-for="(h, i) in headers" :key="i" class="flex gap-2">
            <input v-model="h.key" placeholder="Key" class="w-1/2 bg-surface border border-border rounded-lg px-2 py-1.5 text-xs font-mono focus:outline-none" />
            <input v-model="h.value" placeholder="Value" class="w-1/2 bg-surface border border-border rounded-lg px-2 py-1.5 text-xs font-mono focus:outline-none" />
            <button @click="headers.splice(i, 1)" class="text-muted hover:text-del px-1">✕</button>
          </div>
          <button @click="headers.push({ key: '', value: '' })" class="text-xs text-accent hover:underline font-mono">+ add header</button>
        </div>
      </div>

      <div class="flex flex-col min-h-0">
        <div class="text-xs text-muted uppercase tracking-wide mb-2 font-mono">Body</div>
        <textarea
          v-model="body"
          placeholder="{ }"
          class="flex-1 min-h-[160px] bg-surface border border-border rounded-lg px-3 py-2 text-xs font-mono focus:outline-none focus:ring-2 focus:ring-accent/50 focus:border-accent resize-none"
        ></textarea>
      </div>
    </div>

    <p v-if="errorMsg" class="px-4 pb-2 text-sm text-del">{{ errorMsg }}</p>

    <ResponseViewer v-if="response" :response="response" />
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import client from '../api/client'
import ResponseViewer from './ResponseViewer.vue'

const props = defineProps({
  collections: { type: Array, default: () => [] },
})
const emit = defineEmits(['saved'])

const methods = ['GET', 'POST', 'PUT', 'PATCH', 'DELETE']

const savedRequestId = ref(null)
const collectionId = ref(null)
const name = ref('')
const method = ref('GET')
const url = ref('')
const headers = ref([{ key: '', value: '' }])
const body = ref('')
const response = ref(null)
const sending = ref(false)
const saving = ref(false)
const errorMsg = ref('')

const methodTextClass = computed(() => {
  switch (method.value) {
    case 'GET':
      return 'text-get'
    case 'POST':
      return 'text-post'
    case 'PUT':
      return 'text-put'
    case 'PATCH':
      return 'text-patch'
    case 'DELETE':
      return 'text-del'
    default:
      return ''
  }
})

function headersObject() {
  const obj = {}
  for (const h of headers.value) {
    if (h.key.trim()) obj[h.key.trim()] = h.value
  }
  return obj
}

function headersFromObject(obj) {
  const entries = Object.entries(obj || {})
  return entries.length ? entries.map(([key, value]) => ({ key, value })) : [{ key: '', value: '' }]
}

async function send() {
  errorMsg.value = ''
  response.value = null
  sending.value = true
  try {
    const res = await client.post('/send', {
      method: method.value,
      url: url.value,
      headers: headersObject(),
      body: body.value,
    })
    response.value = res.data.data
  } catch (e) {
    errorMsg.value = e.response?.data?.message || 'Request failed'
  } finally {
    sending.value = false
  }
}

async function save() {
  if (!name.value || !collectionId.value) return
  saving.value = true
  errorMsg.value = ''
  try {
    const payload = {
      collection_id: collectionId.value,
      name: name.value,
      method: method.value,
      url: url.value,
      headers: headersObject(),
      body: body.value,
    }
    let res
    if (savedRequestId.value) {
      res = await client.put(`/requests/${savedRequestId.value}`, payload)
    } else {
      res = await client.post('/requests', payload)
      savedRequestId.value = res.data.data.id
    }
    emit('saved', res.data.data)
  } catch (e) {
    errorMsg.value = e.response?.data?.message || 'Save failed'
  } finally {
    saving.value = false
  }
}

function loadRequest(req) {
  savedRequestId.value = req.id ?? null
  collectionId.value = req.collection_id ?? null
  name.value = req.name ?? ''
  method.value = req.method ?? 'GET'
  url.value = req.url ?? ''
  body.value = req.body ?? ''
  response.value = null
  errorMsg.value = ''
  try {
    const parsed = typeof req.headers === 'string' ? JSON.parse(req.headers || '{}') : req.headers
    headers.value = headersFromObject(parsed)
  } catch {
    headers.value = [{ key: '', value: '' }]
  }
}

function newRequest(presetCollectionId = null) {
  savedRequestId.value = null
  collectionId.value = presetCollectionId
  name.value = ''
  method.value = 'GET'
  url.value = ''
  headers.value = [{ key: '', value: '' }]
  body.value = ''
  response.value = null
  errorMsg.value = ''
}

defineExpose({ loadRequest, newRequest })
</script>
