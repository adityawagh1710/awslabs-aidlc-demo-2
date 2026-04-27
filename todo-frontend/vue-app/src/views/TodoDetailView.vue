<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  ArrowLeftIcon, PencilIcon, TrashIcon, CheckIcon,
  CalendarDaysIcon, TagIcon, PaperClipIcon,
} from '@heroicons/vue/24/outline'
import AppLayout from '@/components/layout/AppLayout.vue'
import AppSpinner from '@/components/ui/AppSpinner.vue'
import PriorityBadge from '@/components/todo/PriorityBadge.vue'
import StatusBadge from '@/components/todo/StatusBadge.vue'
import { todosApi, type Todo, type Status } from '@/api/todos'
import { useTodosStore } from '@/stores/todos'
import { useNotifStore } from '@/stores/notifications'

const route  = useRoute()
const router = useRouter()
const store  = useTodosStore()
const notif  = useNotifStore()

const todo    = ref<Todo | null>(null)
const loading = ref(true)
const editing = ref(false)
const saving  = ref(false)

const editTitle       = ref('')
const editDescription = ref('')
const editPriority    = ref<'low' | 'medium' | 'high'>('medium')
const editDueDate     = ref('')
const editStatus      = ref<Status>('pending')

onMounted(async () => {
  try {
    const { data } = await todosApi.get(route.params.id as string)
    todo.value = data
  } catch { notif.error('Task not found'); router.push('/todos') }
  finally { loading.value = false }
})

function startEdit() {
  if (!todo.value) return
  editTitle.value       = todo.value.title
  editDescription.value = todo.value.description
  editPriority.value    = todo.value.priority
  editDueDate.value     = todo.value.due_date?.substring(0, 10) ?? ''
  editStatus.value      = todo.value.status
  editing.value         = true
}

async function saveEdit() {
  if (!todo.value) return
  saving.value = true
  try {
    const updated = await store.updateTodo(todo.value.id, {
      title:       editTitle.value,
      description: editDescription.value,
      priority:    editPriority.value,
      due_date:    editDueDate.value || undefined,
      status:      editStatus.value,
    })
    todo.value    = updated
    editing.value = false
    notif.success('Task updated')
  } catch { notif.error('Failed to update task') }
  finally { saving.value = false }
}

async function deleteTodo() {
  if (!todo.value) return
  if (!confirm('Delete this task?')) return
  await store.removeTodo(todo.value.id)
  notif.success('Task deleted')
  router.push('/todos')
}

function fmt(d?: string) {
  if (!d) return '—'
  return new Date(d).toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' })
}
</script>

<template>
  <AppLayout>
    <div class="p-8 max-w-2xl mx-auto w-full">
      <!-- Back -->
      <button @click="router.push('/todos')" class="btn-ghost btn-sm mb-6 -ml-2">
        <ArrowLeftIcon class="w-4 h-4" />
        Back to tasks
      </button>

      <!-- Loading -->
      <div v-if="loading" class="flex justify-center py-16"><AppSpinner size="lg" /></div>

      <!-- Content -->
      <div v-else-if="todo" class="card">
        <!-- Header -->
        <div class="px-6 pt-6 pb-4 border-b border-slate-100 flex items-start gap-4">
          <div class="flex-1 min-w-0">
            <div v-if="!editing">
              <h1 class="text-xl font-bold text-slate-800 mb-2" :class="{ 'line-through text-slate-400': todo.status === 'done' }">
                {{ todo.title }}
              </h1>
              <div class="flex items-center flex-wrap gap-2">
                <PriorityBadge :priority="todo.priority" />
                <StatusBadge   :status="todo.status" />
                <span v-if="todo.due_date" class="badge badge-slate">
                  <CalendarDaysIcon class="w-3 h-3" />
                  {{ fmt(todo.due_date) }}
                </span>
              </div>
            </div>
            <!-- Edit mode -->
            <div v-else class="space-y-3">
              <input v-model="editTitle" class="input text-lg font-semibold" />
              <div class="grid grid-cols-3 gap-2">
                <select v-model="editStatus" class="input text-sm">
                  <option value="pending">Pending</option>
                  <option value="in_progress">In Progress</option>
                  <option value="done">Done</option>
                </select>
                <select v-model="editPriority" class="input text-sm">
                  <option value="low">Low</option>
                  <option value="medium">Medium</option>
                  <option value="high">High</option>
                </select>
                <input v-model="editDueDate" type="date" class="input text-sm" />
              </div>
            </div>
          </div>
          <!-- Actions -->
          <div class="flex items-center gap-2 shrink-0">
            <template v-if="editing">
              <button @click="editing = false" class="btn-sm btn-secondary">Cancel</button>
              <button @click="saveEdit" class="btn-sm btn-primary" :disabled="saving">
                <AppSpinner v-if="saving" size="sm" />
                <CheckIcon v-else class="w-4 h-4" />
                Save
              </button>
            </template>
            <template v-else>
              <button @click="startEdit" class="btn-ghost btn-sm"><PencilIcon class="w-4 h-4" /></button>
              <button @click="deleteTodo" class="btn-ghost btn-sm text-red-500 hover:bg-red-50"><TrashIcon class="w-4 h-4" /></button>
            </template>
          </div>
        </div>

        <!-- Body -->
        <div class="px-6 py-5 space-y-5">
          <!-- Description -->
          <div>
            <p class="label">Description</p>
            <div v-if="!editing">
              <p v-if="todo.description" class="text-sm text-slate-600 leading-relaxed">{{ todo.description }}</p>
              <p v-else class="text-sm text-slate-400 italic">No description</p>
            </div>
            <textarea v-else v-model="editDescription" class="input resize-none" rows="3" placeholder="Add details…" />
          </div>

          <!-- Tags -->
          <div v-if="todo.tags.length">
            <p class="label">Tags</p>
            <div class="flex flex-wrap gap-1.5">
              <span v-for="tag in todo.tags" :key="tag.id" class="badge-purple">
                <TagIcon class="w-3 h-3" />
                {{ tag.name }}
              </span>
            </div>
          </div>

          <!-- Metadata -->
          <div class="divider" />
          <div class="grid grid-cols-2 gap-4 text-xs text-slate-400">
            <div><span class="font-medium text-slate-500">Created</span><br />{{ fmt(todo.created_at) }}</div>
            <div><span class="font-medium text-slate-500">Updated</span><br />{{ fmt(todo.updated_at) }}</div>
          </div>

          <!-- File attachments placeholder -->
          <div class="p-4 rounded-xl border-2 border-dashed border-slate-200 flex flex-col items-center gap-2 text-slate-400">
            <PaperClipIcon class="w-6 h-6" />
            <p class="text-xs">File attachments coming soon</p>
          </div>
        </div>
      </div>
    </div>
  </AppLayout>
</template>
