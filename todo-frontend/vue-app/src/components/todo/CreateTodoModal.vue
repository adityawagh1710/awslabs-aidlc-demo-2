<script setup lang="ts">
import { ref, computed } from 'vue'
import AppModal from '@/components/ui/AppModal.vue'
import AppSpinner from '@/components/ui/AppSpinner.vue'
import { useTodosStore } from '@/stores/todos'
import { useNotifStore } from '@/stores/notifications'
import type { Priority } from '@/api/todos'

defineProps<{ show: boolean }>()
const emit = defineEmits<{ close: [] }>()

const store = useTodosStore()
const notif = useNotifStore()

const title       = ref('')
const description = ref('')
const priority    = ref<Priority>('medium')
const dueDate     = ref('')
const selectedTags = ref<string[]>([])
const loading     = ref(false)
const error       = ref('')

const priorities: { value: Priority; label: string }[] = [
  { value: 'low',    label: 'Low' },
  { value: 'medium', label: 'Medium' },
  { value: 'high',   label: 'High' },
]

const canSubmit = computed(() => title.value.trim().length > 0 && !loading.value)

async function submit() {
  if (!canSubmit.value) return
  loading.value = true
  error.value   = ''
  try {
    await store.createTodo({
      title:       title.value.trim(),
      description: description.value.trim() || undefined,
      priority:    priority.value,
      due_date:    dueDate.value || undefined,
      tag_ids:     selectedTags.value.length ? selectedTags.value : undefined,
    })
    notif.success('Task created!')
    reset()
    emit('close')
  } catch {
    error.value = 'Failed to create task. Please try again.'
  } finally {
    loading.value = false
  }
}

function reset() {
  title.value = ''; description.value = ''; priority.value = 'medium'
  dueDate.value = ''; selectedTags.value = []; error.value = ''
}

function toggleTag(id: string) {
  selectedTags.value.includes(id)
    ? selectedTags.value = selectedTags.value.filter(t => t !== id)
    : selectedTags.value.push(id)
}
</script>

<template>
  <AppModal title="New Task" :show="show" @close="emit('close')">
    <form @submit.prevent="submit" class="flex flex-col gap-4">
      <!-- Title -->
      <div>
        <label class="label">Title <span class="text-red-500">*</span></label>
        <input v-model="title" type="text" class="input" placeholder="What needs to be done?" autofocus />
      </div>

      <!-- Description -->
      <div>
        <label class="label">Description</label>
        <textarea v-model="description" class="input resize-none" rows="2" placeholder="Add details..." />
      </div>

      <!-- Priority + Due date -->
      <div class="grid grid-cols-2 gap-3">
        <div>
          <label class="label">Priority</label>
          <select v-model="priority" class="input">
            <option v-for="p in priorities" :key="p.value" :value="p.value">{{ p.label }}</option>
          </select>
        </div>
        <div>
          <label class="label">Due date</label>
          <input v-model="dueDate" type="date" class="input" />
        </div>
      </div>

      <!-- Tags -->
      <div v-if="store.tags.length">
        <label class="label">Tags</label>
        <div class="flex flex-wrap gap-1.5">
          <button
            v-for="tag in store.tags" :key="tag.id" type="button"
            @click="toggleTag(tag.id)"
            :class="[
              'badge transition-all duration-150 cursor-pointer',
              selectedTags.includes(tag.id) ? 'badge-purple ring-2 ring-purple-400' : 'badge-slate hover:badge-indigo'
            ]"
          >
            {{ tag.name }}
          </button>
        </div>
      </div>

      <p v-if="error" class="field-error">{{ error }}</p>
    </form>

    <template #footer>
      <button class="btn-md btn-secondary" @click="emit('close')" type="button">Cancel</button>
      <button class="btn-md btn-primary" @click="submit" :disabled="!canSubmit">
        <AppSpinner v-if="loading" size="sm" />
        <span>Create Task</span>
      </button>
    </template>
  </AppModal>
</template>
