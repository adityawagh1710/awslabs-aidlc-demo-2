<script setup lang="ts">
import { XMarkIcon } from '@heroicons/vue/24/outline'
defineProps<{ title: string; show: boolean }>()
const emit = defineEmits<{ close: [] }>()
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="show" class="fixed inset-0 z-40 flex items-center justify-center p-4">
        <!-- Backdrop -->
        <div class="absolute inset-0 bg-black/40 backdrop-blur-sm" @click="emit('close')" />
        <!-- Panel -->
        <div class="relative w-full max-w-lg bg-white rounded-2xl shadow-float animate-slide-up">
          <div class="flex items-center justify-between px-6 py-4 border-b border-slate-100">
            <h2 class="text-base font-semibold text-slate-800">{{ title }}</h2>
            <button @click="emit('close')" class="btn-ghost btn-sm rounded-lg p-1.5">
              <XMarkIcon class="w-5 h-5" />
            </button>
          </div>
          <div class="px-6 py-5">
            <slot />
          </div>
          <div v-if="$slots.footer" class="px-6 pb-5 flex justify-end gap-2">
            <slot name="footer" />
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.modal-enter-active, .modal-leave-active { transition: opacity 0.2s ease; }
.modal-enter-from, .modal-leave-to { opacity: 0; }
</style>
