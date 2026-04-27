<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import {
  PlusIcon, MagnifyingGlassIcon, XMarkIcon,
  InboxIcon, TagIcon,
} from '@heroicons/vue/24/outline'
import AppLayout from '@/components/layout/AppLayout.vue'
import AppSpinner from '@/components/ui/AppSpinner.vue'
import TodoCard from '@/components/todo/TodoCard.vue'
import CreateTodoModal from '@/components/todo/CreateTodoModal.vue'
import { useTodosStore } from '@/stores/todos'
import { useNotifStore } from '@/stores/notifications'
import type { Status, Priority } from '@/api/todos'

const store = useTodosStore()
const notif = useNotifStore()

const showCreate  = ref(false)
const searchQ     = ref('')
const filterStatus   = ref<Status | ''>('')
const filterPriority = ref<Priority | ''>('')
const showTagMgr  = ref(false)
const newTagName  = ref('')

onMounted(async () => {
  await Promise.all([store.fetchTodos(), store.fetchTags()])
})

const visibleTodos = computed(() => {
  let list = store.todos
  if (searchQ.value) {
    const q = searchQ.value.toLowerCase()
    list = list.filter(t => t.title.toLowerCase().includes(q) || t.description?.toLowerCase().includes(q))
  }
  if (filterStatus.value)   list = list.filter(t => t.status   === filterStatus.value)
  if (filterPriority.value) list = list.filter(t => t.priority === filterPriority.value)
  return list
})

const stats = computed(() => ({
  total:       store.todos.length,
  pending:     store.todos.filter(t => t.status === 'pending').length,
  in_progress: store.todos.filter(t => t.status === 'in_progress').length,
  done:        store.todos.filter(t => t.status === 'done').length,
}))

async function handleComplete(id: string) {
  const todo = store.todos.find(t => t.id === id)
  if (!todo) return
  const next = todo.status === 'done' ? 'pending' : (todo.status === 'pending' ? 'in_progress' : 'done')
  try {
    await store.updateTodo(id, { status: next })
  } catch { notif.error('Failed to update status') }
}

async function handleRemove(id: string) {
  try {
    await store.removeTodo(id)
    notif.success('Task deleted')
  } catch { notif.error('Failed to delete task') }
}

async function addTag() {
  if (!newTagName.value.trim()) return
  try {
    await store.createTag(newTagName.value.trim())
    newTagName.value = ''
    notif.success('Tag created')
  } catch { notif.error('Failed to create tag') }
}

async function removeTag(id: string) {
  try {
    await store.removeTag(id)
    notif.success('Tag deleted')
  } catch { notif.error('Failed to delete tag') }
}

function clearFilters() {
  searchQ.value = ''; filterStatus.value = ''; filterPriority.value = ''
}

const hasFilters = computed(() => searchQ.value || filterStatus.value || filterPriority.value)
</script>

<template>
  <AppLayout>
    <div class="p-8 max-w-4xl mx-auto w-full">

      <!-- Header -->
      <div class="flex items-center justify-between mb-6">
        <div>
          <h1 class="text-2xl font-bold text-slate-800">My Tasks</h1>
          <p class="text-sm text-slate-500 mt-0.5">{{ stats.total }} tasks · {{ stats.done }} done</p>
        </div>
        <div class="flex items-center gap-2">
          <button @click="showTagMgr = !showTagMgr" class="btn-md btn-secondary">
            <TagIcon class="w-4 h-4" />
            Tags
          </button>
          <button @click="showCreate = true" class="btn-md btn-primary">
            <PlusIcon class="w-4 h-4" />
            New Task
          </button>
        </div>
      </div>

      <!-- Stats bar -->
      <div class="grid grid-cols-3 gap-3 mb-6">
        <div v-for="s in [
          { label: 'Pending',     value: stats.pending,     color: 'text-slate-600 bg-slate-100' },
          { label: 'In Progress', value: stats.in_progress, color: 'text-indigo-600 bg-indigo-50' },
          { label: 'Done',        value: stats.done,        color: 'text-green-600  bg-green-50' },
        ]" :key="s.label" class="card p-4 flex items-center gap-3">
          <div :class="['w-10 h-10 rounded-xl flex items-center justify-center text-lg font-bold', s.color]">
            {{ s.value }}
          </div>
          <div>
            <p class="text-xs text-slate-500">{{ s.label }}</p>
          </div>
        </div>
      </div>

      <!-- Search + Filters -->
      <div class="flex items-center gap-2 mb-4">
        <div class="relative flex-1">
          <MagnifyingGlassIcon class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
          <input v-model="searchQ" type="text" class="input pl-9" placeholder="Search tasks…" />
        </div>
        <select v-model="filterStatus" class="input w-36">
          <option value="">All statuses</option>
          <option value="pending">Pending</option>
          <option value="in_progress">In Progress</option>
          <option value="done">Done</option>
        </select>
        <select v-model="filterPriority" class="input w-36">
          <option value="">All priorities</option>
          <option value="low">Low</option>
          <option value="medium">Medium</option>
          <option value="high">High</option>
        </select>
        <button v-if="hasFilters" @click="clearFilters" class="btn-ghost btn-md text-slate-500">
          <XMarkIcon class="w-4 h-4" />
          Clear
        </button>
      </div>

      <!-- Tag manager (collapsible) -->
      <Transition name="slide-down">
        <div v-if="showTagMgr" class="card p-4 mb-4 animate-slide-up">
          <div class="flex items-center gap-2 mb-3">
            <h3 class="text-sm font-semibold text-slate-700">Manage Tags</h3>
          </div>
          <div class="flex flex-wrap gap-2 mb-3">
            <div v-for="tag in store.tags" :key="tag.id"
              class="flex items-center gap-1 badge-purple pr-1">
              <span>{{ tag.name }}</span>
              <button @click="removeTag(tag.id)" class="hover:text-red-500 transition-colors ml-1">
                <XMarkIcon class="w-3 h-3" />
              </button>
            </div>
            <p v-if="!store.tags.length" class="text-xs text-slate-400">No tags yet.</p>
          </div>
          <div class="flex gap-2">
            <input v-model="newTagName" type="text" class="input text-sm py-1.5" placeholder="New tag name…"
              @keyup.enter="addTag" />
            <button @click="addTag" class="btn-sm btn-primary">Add</button>
          </div>
        </div>
      </Transition>

      <!-- Todo list -->
      <div v-if="store.loading" class="flex justify-center py-16">
        <AppSpinner size="lg" />
      </div>

      <div v-else-if="visibleTodos.length" class="card divide-y divide-slate-50">
        <TodoCard
          v-for="todo in visibleTodos" :key="todo.id"
          :todo="todo"
          @complete="handleComplete"
          @remove="handleRemove"
        />
      </div>

      <!-- Empty state -->
      <div v-else class="flex flex-col items-center justify-center py-20 text-center">
        <div class="w-16 h-16 rounded-2xl bg-slate-100 flex items-center justify-center mb-4">
          <InboxIcon class="w-8 h-8 text-slate-400" />
        </div>
        <h3 class="text-base font-semibold text-slate-700 mb-1">
          {{ hasFilters ? 'No matching tasks' : 'No tasks yet' }}
        </h3>
        <p class="text-sm text-slate-400 mb-5">
          {{ hasFilters ? 'Try clearing your filters.' : 'Create your first task to get started.' }}
        </p>
        <button v-if="!hasFilters" @click="showCreate = true" class="btn-md btn-primary">
          <PlusIcon class="w-4 h-4" />
          New Task
        </button>
        <button v-else @click="clearFilters" class="btn-md btn-secondary">Clear filters</button>
      </div>
    </div>
  </AppLayout>

  <CreateTodoModal :show="showCreate" @close="showCreate = false" />
</template>
