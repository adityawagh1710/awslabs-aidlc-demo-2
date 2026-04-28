<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { RouterLink } from 'vue-router'
import { EyeIcon, EyeSlashIcon, CheckCircleIcon } from '@heroicons/vue/24/outline'
import AppSpinner from '@/components/ui/AppSpinner.vue'
import { useAuthStore } from '@/stores/auth'
import { useNotifStore } from '@/stores/notifications'

const auth   = useAuthStore()
const notif  = useNotifStore()
const router = useRouter()

const email     = ref('')
const password  = ref('')
const showPass  = ref(false)
const loading   = ref(false)
const error     = ref('')

async function submit() {
  if (!email.value || !password.value) { error.value = 'Please fill in all fields.'; return }
  loading.value = true; error.value = ''
  try {
    const result = await auth.login(email.value, password.value)
    if (result.mfa_required) { router.push('/mfa'); return }
    notif.success('Welcome back!')
    router.push('/todos')
  } catch (e: any) {
    error.value = e.response?.data?.error === 'unauthorized'
      ? 'Invalid email or password.'
      : e.response?.data?.error === 'too many requests'
      ? 'Too many attempts. Please try again later.'
      : 'Something went wrong. Please try again.'
  } finally { loading.value = false }
}
</script>

<template>
  <div class="min-h-screen flex bg-slate-50">
    <!-- Left panel -->
    <div class="hidden lg:flex lg:w-1/2 bg-gradient-to-br from-primary-600 to-primary-800 flex-col justify-between p-12">
      <div class="flex items-center gap-3">
        <div class="w-9 h-9 bg-white/20 rounded-xl flex items-center justify-center">
          <CheckCircleIcon class="w-5 h-5 text-white" />
        </div>
        <span class="text-white font-bold text-lg tracking-tight">Taskly</span>
      </div>
      <div>
        <h1 class="text-4xl font-bold text-white leading-tight mb-4">
          Stay on top of<br/>everything that<br/>matters.
        </h1>
        <p class="text-primary-200 text-lg leading-relaxed">
          Manage tasks, set reminders, attach files, and collaborate — all in one place.
        </p>
        <!-- Feature list -->
        <ul class="mt-8 space-y-3">
          <li v-for="f in ['Smart task management', 'File attachments', 'Real-time notifications', 'Recurring tasks']"
              :key="f" class="flex items-center gap-3 text-primary-100 text-sm">
            <div class="w-5 h-5 rounded-full bg-white/20 flex items-center justify-center shrink-0">
              <CheckCircleIcon class="w-3 h-3 text-white" />
            </div>
            {{ f }}
          </li>
        </ul>
      </div>
      <p class="text-primary-300 text-sm">&copy; {{ new Date().getFullYear() }} Taskly</p>
    </div>

    <!-- Right panel -->
    <div class="flex-1 flex items-center justify-center p-8">
      <div class="w-full max-w-sm">
        <!-- Mobile logo -->
        <div class="flex items-center gap-2 mb-8 lg:hidden">
          <div class="w-8 h-8 rounded-lg bg-primary-600 flex items-center justify-center">
            <CheckCircleIcon class="w-5 h-5 text-white" />
          </div>
          <span class="font-bold text-slate-800 text-lg">Taskly</span>
        </div>

        <div class="mb-8">
          <h2 class="text-2xl font-bold text-slate-800 mb-1">Sign in</h2>
          <p class="text-slate-500 text-sm">Welcome back! Enter your credentials to continue.</p>
        </div>

        <form @submit.prevent="submit" class="space-y-4">
          <div>
            <label class="label" for="email">Email</label>
            <input id="email" v-model="email" type="email" class="input" placeholder="you@example.com" autocomplete="email" />
          </div>

          <div>
            <label class="label" for="password">Password</label>
            <div class="relative">
              <input id="password" v-model="password" :type="showPass ? 'text' : 'password'"
                class="input pr-10" placeholder="••••••••" autocomplete="current-password" />
              <button type="button" @click="showPass = !showPass"
                class="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 transition-colors">
                <EyeSlashIcon v-if="showPass" class="w-4 h-4" />
                <EyeIcon v-else class="w-4 h-4" />
              </button>
            </div>
          </div>

          <div v-if="error" class="flex items-center gap-2 p-3 rounded-lg bg-red-50 border border-red-100">
            <p class="text-sm text-red-600">{{ error }}</p>
          </div>

          <button type="submit" class="btn-primary btn-lg w-full mt-2" :disabled="loading">
            <AppSpinner v-if="loading" size="sm" />
            <span>{{ loading ? 'Signing in…' : 'Sign in' }}</span>
          </button>
        </form>

        <p class="text-center text-sm text-slate-500 mt-6">
          Don't have an account?
          <RouterLink to="/register" class="text-primary-600 font-medium hover:underline">Sign up</RouterLink>
        </p>
      </div>
    </div>
  </div>
</template>
