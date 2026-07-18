<template>
  <aside class="w-72 bg-surface border-r border-border flex flex-col h-full">
    <div class="p-4 border-b border-border flex items-center justify-between">
      <span class="font-display font-semibold">API Test Studio</span>
    </div>

    <div class="p-3 border-b border-border">
      <button
        @click="showHistory"
        class="w-full text-left text-sm px-3 py-2 rounded-lg hover:bg-surface2 transition-colors flex items-center gap-2"
        :class="{ 'bg-surface2': activeView === 'history' }"
      >
        <span class="font-mono text-muted">⏱</span> History
      </button>
    </div>

    <div class="flex-1 overflow-y-auto p-3">
      <div class="flex items-center justify-between mb-2">
        <span class="text-xs text-muted uppercase tracking-wide font-mono">Collections</span>
        <button @click="creatingCollection = true" class="text-accent text-xs hover:underline">+ New</button>
      </div>

      <div v-if="creatingCollection" class="mb-2 flex gap-1">
        <input
          v-model="newCollectionName"
          @keyup.enter="createCollection"
          placeholder="Collection name"
          class="flex-1 bg-bg border border-border rounded px-2 py-1 text-xs focus:outline-none"
        />
        <button @click="createCollection" class="text-xs text-accent px-1">✓</button>
        <button @click="creatingCollection = false" class="text-xs text-muted px-1">✕</button>
      </div>

      <div v-if="loading" class="text-xs text-muted font-mono">Loading…</div>
      <div v-else-if="!collections.length" class="text-xs text-muted">No collections yet. Create one to start saving requests.</div>

      <div v-for="c in collections" :key="c.id" class="mb-1">
        <div class="group flex items-center justify-between px-2 py-1.5 rounded-lg hover:bg-surface2 cursor-pointer" @click="toggleCollection(c)">
          <span class="text-sm truncate">{{ expanded[c.id] ? '▾' : '▸' }} {{ c.name }}</span>
          <div class="hidden group-hover:flex items-center gap-1">
            <button @click.stop="emit('new-request', c.id)" class="text-xs text-muted hover:text-accent px-1" title="New request">+</button>
            <button @click.stop="deleteCollection(c)" class="text-xs text-muted hover:text-del px-1" title="Delete collection">✕</button>
          </div>
        </div>

        <div v-if="expanded[c.id]" class="ml-4 mt-1 space-y-0.5">
          <div v-if="!requestsByCollection[c.id]?.length" class="text-xs text-muted px-2 py-1">No requests yet.</div>
          <div
            v-for="r in requestsByCollection[c.id]"
            :key="r.id"
            @click="emit('open-request', r)"
            class="flex items-center gap-2 px-2 py-1 rounded hover:bg-surface2 cursor-pointer text-xs"
          >
            <MethodBadge :method="r.method" />
            <span class="truncate text-muted">{{ r.name }}</span>
          </div>
        </div>
      </div>
    </div>

    <div class="p-3 border-t border-border">
      <div class="text-xs text-muted mb-2 truncate">{{ auth.user?.username || auth.user?.email }}</div>
      <button @click="logout" class="text-xs text-muted hover:text-del">Log out</button>
    </div>
  </aside>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import client from '../api/client'
import { useAuthStore } from '../stores/auth'
import MethodBadge from './MethodBadge.vue'

const emit = defineEmits(['open-request', 'new-request', 'show-history', 'collections-loaded'])
defineProps({ activeView: { type: String, default: 'builder' } })

const router = useRouter()
const auth = useAuthStore()

const collections = ref([])
const loading = ref(true)
const expanded = reactive({})
const requestsByCollection = reactive({})
const creatingCollection = ref(false)
const newCollectionName = ref('')

async function loadCollections() {
  loading.value = true
  try {
    const res = await client.get('/collections')
    collections.value = res.data.data || []
    emit('collections-loaded', collections.value)
  } finally {
    loading.value = false
  }
}

async function toggleCollection(c) {
  expanded[c.id] = !expanded[c.id]
  if (expanded[c.id] && !requestsByCollection[c.id]) {
    const res = await client.get(`/requests?collection_id=${c.id}`)
    requestsByCollection[c.id] = res.data.data || []
  }
}

async function createCollection() {
  if (!newCollectionName.value.trim()) return
  const res = await client.post('/collections', { name: newCollectionName.value.trim() })
  collections.value.unshift(res.data.data)
  newCollectionName.value = ''
  creatingCollection.value = false
  emit('collections-loaded', collections.value)
}

async function deleteCollection(c) {
  if (!confirm(`Delete "${c.name}"? This cannot be undone.`)) return
  await client.delete(`/collections/${c.id}`)
  collections.value = collections.value.filter((x) => x.id !== c.id)
  delete requestsByCollection[c.id]
  emit('collections-loaded', collections.value)
}

function showHistory() {
  emit('show-history')
}

function logout() {
  auth.logout()
  router.push({ name: 'login' })
}

function refreshCollection(collectionId) {
  if (expanded[collectionId]) {
    client.get(`/requests?collection_id=${collectionId}`).then((res) => {
      requestsByCollection[collectionId] = res.data.data || []
    })
  }
}

defineExpose({ refreshCollection, loadCollections })

onMounted(loadCollections)
</script>
