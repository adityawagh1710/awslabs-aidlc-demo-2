import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  server: {
    proxy: {
      // Backend API: strip /api so the dev server forwards /api/auth/login
      // → http://localhost:3000/auth/login. Mirrors nginx in the prod image.
      '/api/auth': { target: 'http://localhost:3000', rewrite: p => p.replace(/^\/api/, '') },
      '/api/todos': { target: 'http://localhost:3001', rewrite: p => p.replace(/^\/api/, '') },
      '/api/tags':  { target: 'http://localhost:3001', rewrite: p => p.replace(/^\/api/, '') },
      '/api/files': { target: 'http://localhost:3002', rewrite: p => p.replace(/^\/api/, '') },
      '/ws':        { target: 'ws://localhost:3003', ws: true },
    },
  },
})
