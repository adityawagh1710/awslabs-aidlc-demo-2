import axios from 'axios'

// Namespace API under /api/* so SPA browser routes (/todos, /tags, ...)
// can never collide with backend paths handled by nginx.
const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL ?? '/api',
  timeout: 10_000,
})

api.interceptors.request.use(cfg => {
  const token = sessionStorage.getItem('access_token')
  if (token) cfg.headers.Authorization = `Bearer ${token}`
  return cfg
})

api.interceptors.response.use(
  r => r,
  async err => {
    const original = err.config
    if (err.response?.status === 401 && !original._retry) {
      original._retry = true
      try {
        const refresh = sessionStorage.getItem('refresh_token')
        const { data } = await axios.post(`${api.defaults.baseURL}/auth/refresh`, { refresh_token: refresh })
        sessionStorage.setItem('access_token', data.access_token)
        sessionStorage.setItem('refresh_token', data.refresh_token)
        original.headers.Authorization = `Bearer ${data.access_token}`
        return api(original)
      } catch {
        sessionStorage.clear()
        window.location.href = '/login'
      }
    }
    return Promise.reject(err)
  }
)

export default api
