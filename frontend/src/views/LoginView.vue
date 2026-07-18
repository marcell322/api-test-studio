<template>
  <div class="min-h-screen flex">
    <!-- Brand panel -->
    <div class="hidden lg:flex lg:w-1/2 flex-col justify-between bg-surface border-r border-border p-12">
      <div class="font-display text-xl font-semibold tracking-tight">API Test Studio</div>
      <div>
        <div class="font-mono text-sm text-muted mb-3">// send your first request</div>
        <div class="flex items-center gap-3 bg-bg border border-border rounded-lg px-4 py-3 font-mono text-sm">
          <span class="px-2 py-0.5 rounded bg-get/20 text-get font-semibold">GET</span>
          <span class="text-text">https://api.example.com/users</span>
        </div>
        <div class="mt-4 flex items-center gap-3 text-sm font-mono">
          <span class="px-2 py-0.5 rounded bg-ok/20 text-ok font-semibold">200</span>
          <span class="text-muted">142ms</span>
        </div>
        <p class="mt-8 text-muted text-sm leading-relaxed max-w-sm">
          A lightweight space to build, save, and replay API requests — organized into collections, with every call logged.
        </p>
      </div>
      <div class="text-xs text-muted font-mono">v0.1</div>
    </div>

    <!-- Form panel -->
    <div class="flex-1 flex items-center justify-center p-8">
      <div class="w-full max-w-sm">
        <h1 class="font-display text-2xl font-semibold mb-1">Welcome back</h1>
        <p class="text-muted text-sm mb-8">Log in to continue to your workspace.</p>

        <form @submit.prevent="handleSubmit" class="space-y-4">
          <div>
            <label class="block text-sm text-muted mb-1.5">Email</label>
            <input
              v-model="email"
              type="email"
              required
              class="w-full bg-surface border border-border rounded-lg px-3 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-accent/50 focus:border-accent"
              placeholder="you@example.com"
            />
          </div>
          <div>
            <label class="block text-sm text-muted mb-1.5">Password</label>
            <input
              v-model="password"
              type="password"
              required
              class="w-full bg-surface border border-border rounded-lg px-3 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-accent/50 focus:border-accent"
              placeholder="••••••••"
            />
          </div>

          <p v-if="error" class="text-sm text-del">{{ error }}</p>

          <button
            type="submit"
            :disabled="loading"
            class="w-full bg-accent hover:bg-accent-dim transition-colors rounded-lg px-4 py-2.5 text-sm font-medium text-bg disabled:opacity-50"
          >
            {{ loading ? 'Logging in…' : 'Log in' }}
          </button>
        </form>

        <p class="mt-6 text-sm text-muted">
          Don't have an account?
          <router-link to="/register" class="text-accent hover:underline">Create one</router-link>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const email = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)
const router = useRouter()
const auth = useAuthStore()

async function handleSubmit() {
  error.value = ''
  loading.value = true
  try {
    await auth.login(email.value, password.value)
    router.push({ name: 'workspace' })
  } catch (e) {
    error.value = e.response?.data?.message || 'Login failed'
  } finally {
    loading.value = false
  }
}
</script>
