import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi, type User } from '@/api/auth'

export const useAuthStore = defineStore('auth', () => {
  const user         = ref<User | null>(null)
  const accessToken  = ref(sessionStorage.getItem('access_token'))
  const refreshToken = ref(sessionStorage.getItem('refresh_token'))

  const isAuthenticated = computed(() => !!accessToken.value)

  async function login(email: string, password: string, mfaCode?: string) {
    const { data } = await authApi.login(email, password, mfaCode)
    if (data.mfa_required) return { mfa_required: true }
    setTokens(data.access_token, data.refresh_token)
    user.value = data.user
    return { mfa_required: false }
  }

  async function register(email: string, password: string) {
    const { data } = await authApi.register(email, password)
    setTokens(data.access_token, data.refresh_token)
    user.value = data.user
  }

  async function logout() {
    if (refreshToken.value) {
      try { await authApi.logout(refreshToken.value) } catch { /* ignore */ }
    }
    clearTokens()
  }

  function setTokens(at: string, rt: string) {
    accessToken.value  = at
    refreshToken.value = rt
    sessionStorage.setItem('access_token',  at)
    sessionStorage.setItem('refresh_token', rt)
  }

  function clearTokens() {
    accessToken.value  = null
    refreshToken.value = null
    user.value         = null
    sessionStorage.clear()
  }

  return { user, accessToken, isAuthenticated, login, register, logout }
})
