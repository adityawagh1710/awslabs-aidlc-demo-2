<script setup lang="ts">
import { ref } from 'vue'
import { ShieldCheckIcon, UserCircleIcon, QrCodeIcon, XCircleIcon } from '@heroicons/vue/24/outline'
import QRCode from 'qrcode'
import AppLayout from '@/components/layout/AppLayout.vue'
import AppSpinner from '@/components/ui/AppSpinner.vue'
import { useAuthStore } from '@/stores/auth'
import { useNotifStore } from '@/stores/notifications'
import { authApi } from '@/api/auth'

const auth  = useAuthStore()
const notif = useNotifStore()

const mfaLoading    = ref(false)
const mfaQRData     = ref('')
const mfaSecret     = ref('')
const mfaCode       = ref('')
const disableCode   = ref('')
const showDisable   = ref(false)
const mfaStep       = ref<'idle' | 'scan' | 'done'>(auth.mfaEnabled ? 'done' : 'idle')

async function startMfa() {
  mfaLoading.value = true
  try {
    const { data } = await authApi.enrollMfa()
    mfaSecret.value = data.secret
    mfaQRData.value = await QRCode.toDataURL(data.qr_url, { width: 200, margin: 2 })
    mfaStep.value   = 'scan'
  } catch { notif.error('Could not start MFA setup') }
  finally { mfaLoading.value = false }
}

async function verifyMfa() {
  if (mfaCode.value.length !== 6) return
  mfaLoading.value = true
  try {
    await authApi.verifyMfa(mfaCode.value)
    mfaStep.value = 'done'
    mfaCode.value = ''
    auth.setMfaEnabled(true)
    notif.success('MFA enabled successfully!')
  } catch { notif.error('Invalid code. Please try again.') }
  finally { mfaLoading.value = false }
}

async function disableMfa() {
  if (disableCode.value.length !== 6) return
  mfaLoading.value = true
  try {
    await authApi.disableMfa(disableCode.value)
    mfaStep.value = 'idle'
    disableCode.value = ''
    showDisable.value = false
    auth.setMfaEnabled(false)
    notif.success('MFA disabled.')
  } catch { notif.error('Invalid code. Please try again.') }
  finally { mfaLoading.value = false }
}
</script>

<template>
  <AppLayout>
    <div class="p-8 max-w-2xl mx-auto w-full">
      <h1 class="text-2xl font-bold text-slate-800 mb-6">Settings</h1>

      <!-- Profile section -->
      <section class="card mb-4">
        <div class="px-6 py-4 border-b border-slate-100 flex items-center gap-3">
          <UserCircleIcon class="w-5 h-5 text-slate-400" />
          <h2 class="text-sm font-semibold text-slate-700">Profile</h2>
        </div>
        <div class="px-6 py-5 space-y-3">
          <div>
            <label class="label">Email</label>
            <input :value="auth.email" type="email" class="input" disabled />
          </div>
        </div>
      </section>

      <!-- MFA section -->
      <section class="card">
        <div class="px-6 py-4 border-b border-slate-100 flex items-center gap-3">
          <ShieldCheckIcon class="w-5 h-5 text-slate-400" />
          <h2 class="text-sm font-semibold text-slate-700">Two-factor authentication</h2>
        </div>
        <div class="px-6 py-5">

          <!-- Idle: not enabled -->
          <div v-if="mfaStep === 'idle'">
            <p class="text-sm text-slate-500 mb-4">
              Add an extra layer of security by enabling TOTP-based 2FA.
            </p>
            <button @click="startMfa" class="btn-md btn-secondary" :disabled="mfaLoading">
              <AppSpinner v-if="mfaLoading" size="sm" />
              <QrCodeIcon v-else class="w-4 h-4" />
              Enable 2FA
            </button>
          </div>

          <!-- QR scan step -->
          <div v-else-if="mfaStep === 'scan'" class="space-y-4">
            <p class="text-sm text-slate-600">
              Scan this QR code with your authenticator app (Google Authenticator, Authy, etc.), then enter the 6-digit code below.
            </p>
            <div class="flex justify-center">
              <img v-if="mfaQRData" :src="mfaQRData" alt="MFA QR code" class="w-48 h-48 rounded-xl border border-slate-200 p-2" />
              <div v-else class="w-48 h-48 rounded-xl border border-slate-200 flex items-center justify-center">
                <AppSpinner size="lg" />
              </div>
            </div>
            <p class="text-xs text-slate-400 text-center break-all">Manual key: {{ mfaSecret }}</p>
            <div>
              <label class="label">Verification code</label>
              <input v-model="mfaCode" type="text" inputmode="numeric" maxlength="6"
                class="input text-center text-xl tracking-[0.4em] font-mono" placeholder="000000"
                @keyup.enter="verifyMfa" />
            </div>
            <button @click="verifyMfa" class="btn-md btn-primary w-full" :disabled="mfaLoading || mfaCode.length !== 6">
              <AppSpinner v-if="mfaLoading" size="sm" />
              Verify & activate
            </button>
          </div>

          <!-- Done: MFA enabled -->
          <div v-else class="space-y-4">
            <div class="flex items-center gap-3 p-4 bg-green-50 rounded-xl border border-green-100">
              <ShieldCheckIcon class="w-6 h-6 text-green-600 shrink-0" />
              <div class="flex-1">
                <p class="text-sm font-semibold text-green-800">2FA is enabled</p>
                <p class="text-xs text-green-600">Your account is protected with TOTP authentication.</p>
              </div>
            </div>

            <!-- Disable 2FA -->
            <div v-if="!showDisable">
              <button @click="showDisable = true" class="btn-md btn-ghost text-red-500 hover:bg-red-50">
                <XCircleIcon class="w-4 h-4" />
                Disable 2FA
              </button>
            </div>
            <div v-else class="p-4 rounded-xl border border-red-100 bg-red-50 space-y-3">
              <p class="text-sm text-red-700">Enter your current authenticator code to confirm disabling 2FA.</p>
              <input v-model="disableCode" type="text" inputmode="numeric" maxlength="6"
                class="input text-center text-xl tracking-[0.4em] font-mono" placeholder="000000"
                @keyup.enter="disableMfa" />
              <div class="flex gap-2">
                <button @click="showDisable = false; disableCode = ''" class="btn-md btn-secondary flex-1">Cancel</button>
                <button @click="disableMfa" class="btn-md bg-red-600 text-white hover:bg-red-700 flex-1"
                  :disabled="mfaLoading || disableCode.length !== 6">
                  <AppSpinner v-if="mfaLoading" size="sm" />
                  Confirm disable
                </button>
              </div>
            </div>
          </div>

        </div>
      </section>
    </div>
  </AppLayout>
</template>
