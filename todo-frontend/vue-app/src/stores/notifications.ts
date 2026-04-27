import { defineStore } from 'pinia'
import { ref } from 'vue'

export type NotifType = 'success' | 'error' | 'info'

interface Notif { id: number; type: NotifType; message: string }

export const useNotifStore = defineStore('notif', () => {
  const items = ref<Notif[]>([])
  let seq = 0

  function push(type: NotifType, message: string) {
    const id = ++seq
    items.value.push({ id, type, message })
    setTimeout(() => dismiss(id), 4000)
  }

  function dismiss(id: number) {
    items.value = items.value.filter(n => n.id !== id)
  }

  const success = (m: string) => push('success', m)
  const error   = (m: string) => push('error', m)
  const info    = (m: string) => push('info', m)

  return { items, dismiss, success, error, info }
})
