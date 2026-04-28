<script setup lang="ts">
import { RouterLink, useRouter } from 'vue-router'
import {
  CheckCircleIcon, Cog6ToothIcon, ArrowRightOnRectangleIcon,
} from '@heroicons/vue/24/outline'
import { useAuthStore } from '@/stores/auth'
import { useNotifStore } from '@/stores/notifications'

const auth  = useAuthStore()
const notif = useNotifStore()
const router = useRouter()

async function logout() {
  await auth.logout()
  notif.success('Signed out')
  router.push('/login')
}
</script>

<template>
  <aside class="w-60 shrink-0 flex flex-col bg-white border-r border-slate-100 min-h-screen">
    <!-- Logo -->
    <div class="flex items-center gap-2.5 px-5 py-5 border-b border-slate-100">
      <div class="w-8 h-8 rounded-lg bg-primary-600 flex items-center justify-center">
        <CheckCircleIcon class="w-5 h-5 text-white" />
      </div>
      <span class="text-base font-bold text-slate-800 tracking-tight">Taskly</span>
    </div>

    <!-- Nav -->
    <nav class="flex-1 px-3 py-4 flex flex-col gap-1">
      <RouterLink to="/todos"
        class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-all duration-150"
        :class="$route.path.startsWith('/todos') ? 'bg-primary-50 text-primary-700' : 'text-slate-600 hover:bg-slate-50 hover:text-slate-800'"
      >
        <CheckCircleIcon class="w-5 h-5" />
        My Tasks
      </RouterLink>
    </nav>

    <!-- Bottom actions -->
    <div class="px-3 py-4 border-t border-slate-100 flex flex-col gap-1">
      <RouterLink to="/settings"
        class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium text-slate-600 hover:bg-slate-50 hover:text-slate-800 transition-all duration-150"
      >
        <Cog6ToothIcon class="w-5 h-5" />
        Settings
      </RouterLink>
      <button @click="logout"
        class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium text-slate-600 hover:bg-red-50 hover:text-red-600 transition-all duration-150 w-full text-left"
      >
        <ArrowRightOnRectangleIcon class="w-5 h-5" />
        Sign out
      </button>
    </div>

    <!-- User chip -->
    <div class="px-3 pb-4">
      <div class="flex items-center gap-2.5 px-3 py-2 rounded-xl bg-slate-50">
        <div class="w-7 h-7 rounded-full bg-primary-100 flex items-center justify-center text-primary-700 text-xs font-bold uppercase">
          {{ auth.email?.[0] ?? '?' }}
        </div>
        <p class="text-xs font-medium text-slate-600 truncate">{{ auth.email || 'User' }}</p>
      </div>
    </div>
  </aside>
</template>
