<script setup lang="ts">
import AppSidebar from './AppSidebar.vue'
</script>

<template>
  <div class="flex min-h-screen bg-slate-50">
    <AppSidebar />
    <main class="flex-1 overflow-auto">
      <div class="animate-fade-in">
        <slot />
      </div>
    </main>
  </div>
</template>
