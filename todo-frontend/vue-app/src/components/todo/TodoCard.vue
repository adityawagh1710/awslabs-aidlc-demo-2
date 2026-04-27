<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import { CalendarDaysIcon, TagIcon, TrashIcon, CheckIcon } from '@heroicons/vue/24/outline'
import type { Todo } from '@/api/todos'
import PriorityBadge from './PriorityBadge.vue'
import StatusBadge from './StatusBadge.vue'

const props = defineProps<{ todo: Todo }>()
const emit  = defineEmits<{ complete: [id: string]; remove: [id: string] }>()

const isOverdue = computed(() => {
  if (!props.todo.due_date || props.todo.status === 'done') return false
  return new Date(props.todo.due_date) < new Date()
})

const dueLabel = computed(() => {
  if (!props.todo.due_date) return null
  return new Date(props.todo.due_date).toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
})
</script>

<template>
  <div class="card-hover group p-4 flex items-start gap-3 animate-fade-in">
    <!-- Complete toggle -->
    <button
      @click.stop="emit('complete', todo.id)"
      :class="[
        'w-5 h-5 rounded-full border-2 shrink-0 mt-0.5 flex items-center justify-center transition-all duration-150',
        todo.status === 'done'
          ? 'bg-primary-600 border-primary-600'
          : 'border-slate-300 hover:border-primary-400 hover:bg-primary-50'
      ]"
      :aria-label="todo.status === 'done' ? 'Mark incomplete' : 'Mark complete'"
    >
      <CheckIcon v-if="todo.status === 'done'" class="w-3 h-3 text-white" />
    </button>

    <!-- Content -->
    <RouterLink :to="`/todos/${todo.id}`" class="flex-1 min-w-0">
      <p :class="['text-sm font-medium truncate mb-1', todo.status === 'done' ? 'text-slate-400 line-through' : 'text-slate-800']">
        {{ todo.title }}
      </p>
      <p v-if="todo.description" class="text-xs text-slate-400 truncate mb-2">{{ todo.description }}</p>

      <div class="flex items-center flex-wrap gap-1.5">
        <PriorityBadge :priority="todo.priority" />
        <StatusBadge   :status="todo.status" />

        <span v-if="dueLabel" :class="['badge', isOverdue ? 'badge-red' : 'badge-slate']">
          <CalendarDaysIcon class="w-3 h-3" />
          {{ dueLabel }}
          <span v-if="isOverdue" class="font-semibold">overdue</span>
        </span>

        <span v-for="tag in todo.tags" :key="tag.id" class="badge-purple">
          <TagIcon class="w-3 h-3" />
          {{ tag.name }}
        </span>
      </div>
    </RouterLink>

    <!-- Delete -->
    <button
      @click.stop="emit('remove', todo.id)"
      class="opacity-0 group-hover:opacity-100 btn-ghost btn-sm rounded-lg p-1.5 text-slate-400 hover:text-red-500 transition-all duration-150"
      aria-label="Delete todo"
    >
      <TrashIcon class="w-4 h-4" />
    </button>
  </div>
</template>
