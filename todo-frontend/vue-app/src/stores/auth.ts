import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '@/api/auth'

export const useAuthStore = defineStore('auth', () => {
  const email        = ref(sessionStorage.getItem('user_email') || '')
  const accessToken  = ref(sessionStorage.getItem('access_token'))
  const refreshToken = ref(sessionStorage.getItem('refresh_token'))
  const mfaEnabled   = ref(sessionStorage.getItem('mfa_enabled') === 'true')

  // Temporarily hold credentials while waiting for MFA code
  const pendingEmail    = ref('')
  const pendingPassword = ref('')

  const isAuthenticated = computed(() => !!accessToken.value)

  async function login(emailVal: string, password: string, mfaCode?: string) {
    try {
      const { data } = await authApi.login(emailVal, password, mfaCode)
      setTokens(data.access_token, data.refresh_token)
      setEmail(emailVal)
      if (mfaCode) setMfaEnabled(true)
      pendingEmail.value = ''
      pendingPassword.value = ''
      return { mfa_required: false }
    } catch (err: any) {
      // Backend returns 422 with {"error":"mfa_required"} when MFA is needed
      if (err.response?.status === 422 && err.response?.data?.error === 'mfa_required') {
        pendingEmail.value = emailVal
        pendingPassword.value = password
        return { mfa_required: true }
      }
      throw err
    }
  }

  async function completeMfa(mfaCode: string) {
    if (!pendingEmail.value || !pendingPassword.value) {
      throw new Error('No pending login — go back to the login page.')
    }
    const { data } = await authApi.login(pendingEmail.value, pendingPassword.value, mfaCode)
    setTokens(data.access_token, data.refresh_token)
    setEmail(pendingEmail.value)
    setMfaEnabled(true)
    pendingEmail.value = ''
    pendingPassword.value = ''
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
    pendingEmail.value = ''
    pendingPassword.value = ''
    sessionStorage.clear()
  }

  return {
    email, accessToken, refreshToken, isAuthenticated, mfaEnabled,
    pendingEmail, pendingPassword,
    login, completeMfa, register, logout, setMfaEnabled,
  }
})
