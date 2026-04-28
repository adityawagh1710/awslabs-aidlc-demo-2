<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ShieldCheckIcon } from '@heroicons/vue/24/outline'
import AppSpinner from '@/components/ui/AppSpinner.vue'
import { useAuthStore } from '@/stores/auth'
import { useNotifStore } from '@/stores/notifications'

const auth   = useAuthStore()
const notif  = useNotifStore()
const router = useRouter()

const code    = ref('')
const loading = ref(false)
const error   = ref('')

async function submit() {
  if (code.value.length !== 6) { error.value = 'Enter the 6-digit code from your authenticator app.'; return }
  if (!auth.pendingEmail) { error.value = 'Session expired. Please go back and log in again.'; return }
  loading.value = true; error.value = ''
  try {
    await auth.completeMfa(code.value)
    notif.success('Authenticated!')
    router.push('/todos')
  } catch {
    error.value = 'Invalid code. Please try again.'
  } finally { loading.value = false }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-slate-50 p-4">
    <div class="w-full max-w-sm">
      <div class="card p-8 text-center">
        <div class="w-14 h-14 rounded-2xl bg-primary-50 flex items-center justify-center mx-auto mb-5">
          <ShieldCheckIcon class="w-7 h-7 text-primary-600" />
        </div>
        <h2 class="text-xl font-bold text-slate-800 mb-1">Two-factor authentication</h2>
        <p class="text-sm text-slate-500 mb-6">Enter the 6-digit code from your authenticator app.</p>

        <form @submit.prevent="submit" class="space-y-4">
          <input v-model="code" type="text" inputmode="numeric" maxlength="6"
            class="input text-center text-2xl tracking-[0.5em] font-mono" placeholder="000000" />

          <div v-if="error" class="p-3 rounded-lg bg-red-50 border border-red-100">
            <p class="text-sm text-red-600">{{ error }}</p>
          </div>

          <button type="submit" class="btn-primary btn-lg w-full" :disabled="loading || code.length !== 6">
            <AppSpinner v-if="loading" size="sm" />
            <span>Verify</span>
          </button>
        </form>

        <button @click="router.push('/login')" class="text-sm text-slate-500 hover:text-slate-700 mt-4 mx-auto block">
          ← Back to login
        </button>
      </div>
    </div>
  </div>
</template>
