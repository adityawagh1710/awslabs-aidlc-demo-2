import api from './client'

export type Priority = 'low' | 'medium' | 'high'
export type Status   = 'pending' | 'in_progress' | 'done'

export interface Tag  { id: string; name: string }
export interface Todo {
  id: string; user_id: string; title: string; description: string
  status: Status; priority: Priority; due_date?: string
  tags: Tag[]; created_at: string; updated_at: string
}

export interface CreateTodoInput {
  title: string; description?: string; priority?: Priority
  due_date?: string; tag_ids?: string[]
}
export interface UpdateTodoInput extends Partial<CreateTodoInput> {
  status?: Status
}
export interface TodoFilter {
  status?: Status; priority?: Priority; tag_id?: string
}

// Backend unmarshals due_date into Go's time.Time, which requires RFC3339.
// <input type="date"> produces a YYYY-MM-DD string; widen it to midnight UTC.
function normalizeBody<T extends { due_date?: string }>(b: T): T {
  if (b.due_date && /^\d{4}-\d{2}-\d{2}$/.test(b.due_date)) {
    return { ...b, due_date: `${b.due_date}T00:00:00Z` }
  }
  return b
}

export const todosApi = {
  list:   (f?: TodoFilter) => api.get<Todo[]>('/todos', { params: f }),
  get:    (id: string)     => api.get<Todo>(`/todos/${id}`),
  create: (b: CreateTodoInput)           => api.post<Todo>('/todos', normalizeBody(b)),
  update: (id: string, b: UpdateTodoInput) => api.put<Todo>(`/todos/${id}`, normalizeBody(b)),
  remove: (id: string)     => api.delete(`/todos/${id}`),
  search: (q: string)      => api.get<Todo[]>('/todos/search', { params: { q } }),
}

export const tagsApi = {
  list:   ()              => api.get<Tag[]>('/tags'),
  create: (name: string)  => api.post<Tag>('/tags', { name }),
  remove: (id: string)    => api.delete(`/tags/${id}`),
}
