import { defineStore } from 'pinia'
import client from '../api/client'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('ats_token') || null,
    user: JSON.parse(localStorage.getItem('ats_user') || 'null'),
  }),
  getters: {
    isAuthenticated: (state) => !!state.token,
  },
  actions: {
    async login(email, password) {
      const res = await client.post('/login', { email, password })
      const { token, user } = res.data.data
      this.token = token
      this.user = user
      localStorage.setItem('ats_token', token)
      localStorage.setItem('ats_user', JSON.stringify(user))
    },
    async register(username, email, password) {
      await client.post('/register', { username, email, password })
    },
    logout() {
      this.token = null
      this.user = null
      localStorage.removeItem('ats_token')
      localStorage.removeItem('ats_user')
    },
  },
})
