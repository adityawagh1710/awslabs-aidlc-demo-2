<script setup lang="ts">
import type { Status } from '@/api/todos'
defineProps<{ status: Status }>()
const map = {
  pending:     { cls: 'badge-slate',  label: 'Pending' },
  in_progress: { cls: 'badge-indigo', label: 'In Progress' },
  done:        { cls: 'badge-green',  label: 'Done' },
}
</script>
<template>
  <span :class="map[status].cls">{{ map[status].label }}</span>
</template>
