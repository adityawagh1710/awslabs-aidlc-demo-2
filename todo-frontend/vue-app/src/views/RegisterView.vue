<script setup lang="ts">
import { ref, computed } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { CheckCircleIcon, EyeIcon, EyeSlashIcon } from '@heroicons/vue/24/outline'
import AppSpinner from '@/components/ui/AppSpinner.vue'
import { useAuthStore } from '@/stores/auth'
import { useNotifStore } from '@/stores/notifications'

const auth   = useAuthStore()
const notif  = useNotifStore()
const router = useRouter()

const email     = ref('')
const password  = ref('')
const confirm   = ref('')
const showPass  = ref(false)
const loading   = ref(false)
const error     = ref('')

const strength = computed(() => {
  const p = password.value
  if (!p) return 0
  let s = 0
  if (p.length >= 8)          s++
  if (/[A-Z]/.test(p))        s++
  if (/[0-9]/.test(p))        s++
  if (/[^A-Za-z0-9]/.test(p)) s++
  return s
})
const strengthLabel = ['', 'Weak', 'Fair', 'Good', 'Strong']
const strengthColor = ['', 'bg-red-400', 'bg-amber-400', 'bg-blue-400', 'bg-green-500']

async function submit() {
  error.value = ''
  if (!email.value || !password.value) { error.value = 'Please fill in all fields.'; return }
  if (password.value !== confirm.value) { error.value = 'Passwords do not match.'; return }
  if (password.value.length < 8) { error.value = 'Password must be at least 8 characters.'; return }
  loading.value = true
  try {
    await auth.register(email.value, password.value)
    notif.success('Account created! Welcome to Taskly.')
    router.push('/todos')
  } catch (e: any) {
    error.value = e.response?.data?.error === 'email already registered'
      ? 'This email is already registered. Try signing in instead.'
      : 'Registration failed. Please try again.'
  } finally { loading.value = false }
}
</script>

<template>
  <div class="min-h-screen flex bg-slate-50">
    <!-- Left panel -->
    <div class="hidden lg:flex lg:w-1/2 bg-gradient-to-br from-violet-600 to-primary-700 flex-col justify-between p-12">
      <div class="flex items-center gap-3">
        <div class="w-9 h-9 bg-white/20 rounded-xl flex items-center justify-center">
          <CheckCircleIcon class="w-5 h-5 text-white" />
        </div>
        <span class="text-white font-bold text-lg tracking-tight">Taskly</span>
      </div>
      <div>
        <h1 class="text-4xl font-bold text-white leading-tight mb-4">Start getting<br/>things done<br/>today.</h1>
        <p class="text-violet-200 text-lg leading-relaxed">
          Create your free account and start organizing your life in minutes.
        </p>
      </div>
      <p class="text-violet-300 text-sm">&copy; {{ new Date().getFullYear() }} Taskly</p>
    </div>

    <!-- Right panel -->
    <div class="flex-1 flex items-center justify-center p-8">
      <div class="w-full max-w-sm">
        <div class="mb-8">
          <h2 class="text-2xl font-bold text-slate-800 mb-1">Create account</h2>
          <p class="text-slate-500 text-sm">Free forever. No credit card required.</p>
        </div>

        <form @submit.prevent="submit" class="space-y-4">
          <div>
            <label class="label" for="reg-email">Email</label>
            <input id="reg-email" v-model="email" type="email" class="input" placeholder="you@example.com" autocomplete="email" />
          </div>

          <div>
            <label class="label" for="reg-pass">Password</label>
            <div class="relative">
              <input id="reg-pass" v-model="password" :type="showPass ? 'text' : 'password'"
                class="input pr-10" placeholder="Min. 8 characters" autocomplete="new-password" />
              <button type="button" @click="showPass = !showPass"
                class="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 transition-colors">
                <EyeSlashIcon v-if="showPass" class="w-4 h-4" />
                <EyeIcon v-else class="w-4 h-4" />
              </button>
            </div>
            <!-- Strength bar -->
            <div v-if="password" class="mt-2 flex items-center gap-2">
              <div class="flex gap-1 flex-1">
                <div v-for="i in 4" :key="i"
                  :class="['h-1 flex-1 rounded-full transition-all duration-300', i <= strength ? strengthColor[strength] : 'bg-slate-200']" />
              </div>
              <span class="text-xs text-slate-500">{{ strengthLabel[strength] }}</span>
            </div>
          </div>

          <div>
            <label class="label" for="reg-confirm">Confirm password</label>
            <input id="reg-confirm" v-model="confirm" type="password" class="input" placeholder="Repeat password" autocomplete="new-password" />
          </div>

          <div v-if="error" class="flex items-center gap-2 p-3 rounded-lg bg-red-50 border border-red-100">
            <p class="text-sm text-red-600">{{ error }}</p>
          </div>

          <button type="submit" class="btn-primary btn-lg w-full mt-2" :disabled="loading">
            <AppSpinner v-if="loading" size="sm" />
            <span>{{ loading ? 'Creating account…' : 'Create account' }}</span>
          </button>
        </form>

        <p class="text-center text-sm text-slate-500 mt-6">
          Already have an account?
          <RouterLink to="/login" class="text-primary-600 font-medium hover:underline">Sign in</RouterLink>
        </p>
      </div>
    </div>
  </div>
</template>
