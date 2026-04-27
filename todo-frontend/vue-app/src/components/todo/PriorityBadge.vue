<script setup lang="ts">
import type { Priority } from '@/api/todos'
const props = defineProps<{ priority: Priority }>()
const map = {
  low:    { cls: 'badge-green',  label: 'Low' },
  medium: { cls: 'badge-amber',  label: 'Medium' },
  high:   { cls: 'badge-red',    label: 'High' },
}
</script>
<template>
  <span :class="map[priority].cls">{{ map[priority].label }}</span>
</template>
