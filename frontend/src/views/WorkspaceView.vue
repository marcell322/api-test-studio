<template>
  <div class="flex h-screen">
    <AppSidebar
      ref="sidebarRef"
      :active-view="activeView"
      @open-request="openRequest"
      @new-request="onNewRequest"
      @show-history="activeView = 'history'"
      @collections-loaded="collections = $event"
    />

    <main class="flex-1 flex flex-col min-w-0">
      <RequestBuilder v-show="activeView === 'builder'" ref="builderRef" :collections="collections" @saved="onSaved" />
      <HistoryPanel v-if="activeView === 'history'" ref="historyRef" @replay="onReplay" />
    </main>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import AppSidebar from '../components/AppSidebar.vue'
import RequestBuilder from '../components/RequestBuilder.vue'
import HistoryPanel from '../components/HistoryPanel.vue'

const activeView = ref('builder')
const collections = ref([])
const sidebarRef = ref(null)
const builderRef = ref(null)

function openRequest(req) {
  activeView.value = 'builder'
  builderRef.value?.loadRequest(req)
}

function onNewRequest(collectionId) {
  activeView.value = 'builder'
  builderRef.value?.newRequest(collectionId)
}

function onSaved(saved) {
  sidebarRef.value?.refreshCollection(saved.collection_id)
}

function onReplay(historyItem) {
  activeView.value = 'builder'
  builderRef.value?.loadRequest({
    method: historyItem.method,
    url: historyItem.url,
    headers: {},
    body: '',
    name: '',
    collection_id: null,
    id: null,
  })
}
</script>
