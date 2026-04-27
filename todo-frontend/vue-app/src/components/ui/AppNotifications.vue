<script setup lang="ts">
import { useNotifStore } from '@/stores/notifications'
import { CheckCircleIcon, XCircleIcon, InformationCircleIcon, XMarkIcon } from '@heroicons/vue/24/solid'
const store = useNotifStore()

const icons = { success: CheckCircleIcon, error: XCircleIcon, info: InformationCircleIcon }
const colors = {
  success: 'bg-green-50 border-green-200 text-green-800',
  error:   'bg-red-50 border-red-200 text-red-800',
  info:    'bg-blue-50 border-blue-200 text-blue-800',
}
const iconColors = {
  success: 'text-green-500', error: 'text-red-500', info: 'text-blue-500',
}
</script>

<template>
  <Teleport to="body">
    <div class="fixed top-4 right-4 z-50 flex flex-col gap-2 w-80">
      <TransitionGroup name="notif">
        <div
          v-for="n in store.items" :key="n.id"
          :class="['flex items-start gap-3 p-3.5 rounded-xl border shadow-float animate-slide-up', colors[n.type]]"
        >
          <component :is="icons[n.type]" :class="['w-5 h-5 shrink-0 mt-0.5', iconColors[n.type]]" />
          <p class="text-sm font-medium flex-1">{{ n.message }}</p>
          <button @click="store.dismiss(n.id)" class="text-current opacity-50 hover:opacity-80 transition-opacity">
            <XMarkIcon class="w-4 h-4" />
          </button>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<style scoped>
.notif-enter-active { transition: all 0.25s ease-out; }
.notif-leave-active { transition: all 0.2s ease-in; }
.notif-enter-from  { opacity: 0; transform: translateX(100%); }
.notif-leave-to    { opacity: 0; transform: translateX(100%); }
</style>
