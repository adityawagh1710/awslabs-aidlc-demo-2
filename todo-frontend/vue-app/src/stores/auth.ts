import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '@/api/auth'

export const useAuthStore = defineStore('auth', () => {
  const email        = ref(sessionStorage.getItem('user_email') || '')
  const accessToken  = ref(sessionStorage.getItem('access_token'))
  const refreshToken = ref(sessionStorage.getItem('refresh_token'))
  const mfaEnabled   = ref(sessionStorage.getItem('mfa_enabled') === 'true')

  const isAuthenticated = computed(() => !!accessToken.value)

  async function login(emailVal: string, password: string, mfaCode?: string) {
    const { data } = await authApi.login(emailVal, password, mfaCode)
    if (data.mfa_required) return { mfa_required: true }
    setTokens(data.access_token, data.refresh_token)
    setEmail(emailVal)
    return { mfa_required: false }
  }

  async function register(emailVal: string, password: string) {
    const { data } = await authApi.register(emailVal, password)
    setTokens(data.access_token, data.refresh_token)
    setEmail(emailVal)
  }

  async function logout() {
    if (refreshToken.value) {
      try { await authApi.logout(refreshToken.value) } catch { /* ignore */ }
    }
    clearTokens()
  }

  function setEmail(e: string) {
    email.value = e
    sessionStorage.setItem('user_email', e)
  }

  function setMfaEnabled(v: boolean) {
    mfaEnabled.value = v
    sessionStorage.setItem('mfa_enabled', String(v))
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
    email.value        = ''
    mfaEnabled.value   = false
    sessionStorage.clear()
  }

  return { email, accessToken, refreshToken, isAuthenticated, mfaEnabled, login, register, logout, setMfaEnabled }
})
