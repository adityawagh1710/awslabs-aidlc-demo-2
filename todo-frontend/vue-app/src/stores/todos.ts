import { defineStore } from 'pinia'
import { ref } from 'vue'
import { todosApi, tagsApi, type Todo, type Tag, type CreateTodoInput, type UpdateTodoInput, type TodoFilter } from '@/api/todos'

export const useTodosStore = defineStore('todos', () => {
  const todos       = ref<Todo[]>([])
  const tags        = ref<Tag[]>([])
  const loading     = ref(false)
  const searchQuery = ref('')

  async function fetchTodos(filter?: TodoFilter) {
    loading.value = true
    try {
      const { data } = await todosApi.list(filter)
      todos.value = data
    } finally { loading.value = false }
  }

  async function fetchTags() {
    const { data } = await tagsApi.list()
    tags.value = data
  }

  async function createTodo(input: CreateTodoInput) {
    const { data } = await todosApi.create(input)
    todos.value.unshift(data)
    return data
  }

  async function updateTodo(id: string, input: UpdateTodoInput) {
    const { data } = await todosApi.update(id, input)
    const idx = todos.value.findIndex(t => t.id === id)
    if (idx !== -1) todos.value[idx] = data
    return data
  }

  async function removeTodo(id: string) {
    await todosApi.remove(id)
    todos.value = todos.value.filter(t => t.id !== id)
  }

  async function searchTodos(q: string) {
    loading.value = true
    searchQuery.value = q
    try {
      const { data } = await todosApi.search(q)
      todos.value = data
    } finally { loading.value = false }
  }

  async function createTag(name: string) {
    const { data } = await tagsApi.create(name)
    tags.value.push(data)
    return data
  }

  async function removeTag(id: string) {
    await tagsApi.remove(id)
    tags.value = tags.value.filter(t => t.id !== id)
  }

  return {
    todos, tags, loading, searchQuery,
    fetchTodos, fetchTags, createTodo, updateTodo, removeTodo,
    searchTodos, createTag, removeTag,
  }
})
