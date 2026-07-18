<template>
  <div class="min-h-screen flex items-center justify-center p-8">
    <div class="w-full max-w-sm">
      <h1 class="font-display text-2xl font-semibold mb-1">Create your account</h1>
      <p class="text-muted text-sm mb-8">Start building and saving API requests.</p>

      <form @submit.prevent="handleSubmit" class="space-y-4">
        <div>
          <label class="block text-sm text-muted mb-1.5">Username</label>
          <input
            v-model="username"
            required
            class="w-full bg-surface border border-border rounded-lg px-3 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-accent/50 focus:border-accent"
            placeholder="alice"
          />
        </div>
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
            minlength="6"
            class="w-full bg-surface border border-border rounded-lg px-3 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-accent/50 focus:border-accent"
            placeholder="At least 6 characters"
          />
        </div>

        <p v-if="error" class="text-sm text-del">{{ error }}</p>
        <p v-if="success" class="text-sm text-ok">Account created. You can log in now.</p>

        <button
          type="submit"
          :disabled="loading"
          class="w-full bg-accent hover:bg-accent-dim transition-colors rounded-lg px-4 py-2.5 text-sm font-medium text-bg disabled:opacity-50"
        >
          {{ loading ? 'Creating…' : 'Create account' }}
        </button>
      </form>

      <p class="mt-6 text-sm text-muted">
        Already have an account?
        <router-link to="/login" class="text-accent hover:underline">Log in</router-link>
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const username = ref('')
const email = ref('')
const password = ref('')
const error = ref('')
const success = ref(false)
const loading = ref(false)
const router = useRouter()
const auth = useAuthStore()

async function handleSubmit() {
  error.value = ''
  success.value = false
  loading.value = true
  try {
    await auth.register(username.value, email.value, password.value)
    success.value = true
    setTimeout(() => router.push({ name: 'login' }), 900)
  } catch (e) {
    error.value = e.response?.data?.message || 'Registration failed'
  } finally {
    loading.value = false
  }
}
</script>
